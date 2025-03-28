package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/stjudewashere/seonaut/internal/models"
	"github.com/stjudewashere/seonaut/internal/services"
)

type projectHandler struct {
	*services.Container
}

// indexHandler Handles the user homepage request and lists all the user's projects.
func (h *projectHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	// If the user is crawling or deleting projects, set a meta refresh tag
	// in the HTML so the page is updated with new data every few seconds.
	refresh := h.ProjectViewService.UserIsProcessingProjects(user.Id)

	v := &PageView{
		Data: struct {
			Projects []models.ProjectView
		}{
			Projects: h.ProjectViewService.GetProjectViews(user.Id),
		},
		User:      *user,
		PageTitle: "PROJECTS_VIEW",
		Refresh:   refresh,
	}

	h.Renderer.RenderTemplate(w, "home", v)
}

// addGetHandler displays the form for adding a new project.
// This handler handles the GET request.
func (h *projectHandler) addGetHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pageView := &PageView{
		User:      *user,
		PageTitle: "ADD_PROJECT",
		Data:      &struct{ Error bool }{},
	}

	h.Renderer.RenderTemplate(w, "project_add", pageView)
}

// addPostHandler handles the POST request to add a project.
// If an there's an error updating the project a variable is set to display an error
// message in the HTML template.
func (h *projectHandler) addPostHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("serveProjectAdd ParseForm: %v\n", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ignoreRobotsTxt, err := strconv.ParseBool(r.FormValue("ignore_robotstxt"))
	if err != nil {
		ignoreRobotsTxt = false
	}

	followNofollow, err := strconv.ParseBool(r.FormValue("follow_nofollow"))
	if err != nil {
		followNofollow = false
	}

	includeNoindex, err := strconv.ParseBool(r.FormValue("include_noindex"))
	if err != nil {
		includeNoindex = false
	}

	crawlSitemap, err := strconv.ParseBool(r.FormValue("crawl_sitemap"))
	if err != nil {
		crawlSitemap = false
	}

	allowSubdomains, err := strconv.ParseBool(r.FormValue("allow_subdomains"))
	if err != nil {
		allowSubdomains = false
	}

	checkExternalLinks, err := strconv.ParseBool(r.FormValue("check_external_links"))
	if err != nil {
		checkExternalLinks = false
	}

	basicAuth, err := strconv.ParseBool(r.FormValue("basic_auth"))
	if err != nil {
		basicAuth = false
	}

	archive, err := strconv.ParseBool(r.FormValue("archive"))
	if err != nil {
		archive = false
	}

	project := &models.Project{
		URL:                r.FormValue("url"),
		IgnoreRobotsTxt:    ignoreRobotsTxt,
		FollowNofollow:     followNofollow,
		IncludeNoindex:     includeNoindex,
		CrawlSitemap:       crawlSitemap,
		AllowSubdomains:    allowSubdomains,
		BasicAuth:          basicAuth,
		CheckExternalLinks: checkExternalLinks,
		Archive:            archive,
	}

	err = h.ProjectService.SaveProject(project, user.Id)
	if err != nil {
		pageView := &PageView{
			User:      *user,
			PageTitle: "ADD_PROJECT",
			Data:      &struct{ Error bool }{Error: true},
		}
		h.Renderer.RenderTemplate(w, "project_add", pageView)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// deleteGetHandler handles the deletion of a project.
// It expects a query parameter "pid" containing the project id to be deleted.
func (h *projectHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p, err := h.ProjectService.FindProject(pid, user.Id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.ProjectService.DeleteProject(&p)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// editGetHandler displays the edition form of a project.
// It expects a query parameter "pid" containing the project id to be edited.
// The form is pre-populated with the project's data.
// Thes handler handles the GET requests.
func (h *projectHandler) editGetHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p, err := h.ProjectService.FindProject(pid, user.Id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := &struct {
		Project models.Project
		Error   bool
	}{
		Project: p,
	}

	pageView := &PageView{
		User:      *user,
		PageTitle: "EDIT_PROJECT",
		Data:      data,
	}

	h.Renderer.RenderTemplate(w, "project_edit", pageView)
}

// editPostHandler handles project edits.
// A variable is set if there's an error updating the project, so an error message
// can be displayed in the HTML page.
// This handler handles the POST request.
func (h *projectHandler) editPostHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pid, err := strconv.Atoi(r.URL.Query().Get("pid"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p, err := h.ProjectService.FindProject(pid, user.Id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p.IgnoreRobotsTxt, err = strconv.ParseBool(r.FormValue("ignore_robotstxt"))
	if err != nil {
		p.IgnoreRobotsTxt = false
	}

	p.FollowNofollow, err = strconv.ParseBool(r.FormValue("follow_nofollow"))
	if err != nil {
		p.FollowNofollow = false
	}

	p.IncludeNoindex, err = strconv.ParseBool(r.FormValue("include_noindex"))
	if err != nil {
		p.IncludeNoindex = false
	}

	p.CrawlSitemap, err = strconv.ParseBool(r.FormValue("crawl_sitemap"))
	if err != nil {
		p.CrawlSitemap = false
	}

	p.AllowSubdomains, err = strconv.ParseBool(r.FormValue("allow_subdomains"))
	if err != nil {
		p.AllowSubdomains = false
	}

	p.CheckExternalLinks, err = strconv.ParseBool(r.FormValue("check_external_links"))
	if err != nil {
		p.CheckExternalLinks = false
	}

	p.BasicAuth, err = strconv.ParseBool(r.FormValue("basic_auth"))
	if err != nil {
		p.BasicAuth = false
	}

	p.Archive, err = strconv.ParseBool(r.FormValue("archive"))
	if err != nil {
		p.Archive = false
	}

	err = h.ProjectService.UpdateProject(&p)
	if err != nil {
		pageView := &PageView{
			User:      *user,
			PageTitle: "EDIT_PROJECT",
			Data: &struct {
				Project models.Project
				Error   bool
			}{
				Project: p,
				Error:   true,
			},
		}

		h.Renderer.RenderTemplate(w, "project_edit", pageView)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
