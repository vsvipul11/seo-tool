package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/stjudewashere/seonaut/internal/models"
	"github.com/stjudewashere/seonaut/internal/services"
)

type resourceHandler struct {
	*services.Container
}

// indexHandler handles the HTTP request for the resources view page.
// It expects the following query parameters:
// - "pid" containing the project id.
// - "rid" the id of the resource to be loaded.
// - "eid" the id of the issue type from wich the user loaded this resource.
// - "ep" the explorer page number from which the user loaded this resource.
// - "t" the tab to be loaded, which defaults to the details tab.
// - "p" the number of page to be loaded, in case the resource page has pagination.
func (h *resourceHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
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

	rid, err := strconv.Atoi(r.URL.Query().Get("rid"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	eid := r.URL.Query().Get("eid")
	ep := r.URL.Query().Get("ep")
	if eid == "" && ep == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tab := r.URL.Query().Get("t")
	if tab == "" {
		tab = "details"
	}

	page, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		page = 1
	}

	pv, err := h.ProjectViewService.GetProjectView(pid, user.Id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pageReportView := h.ReportService.GetPageReport(rid, pv.Crawl.Id, tab, page)
	isArchived := h.Container.ArchiveService.ArchiveExists(&pv.Project)
	isTextMedia := strings.HasPrefix(pageReportView.PageReport.MediaType, "text/")

	data := &struct {
		PageReportView *models.PageReportView
		ProjectView    *models.ProjectView
		Eid            string
		Ep             string
		Tab            string
		Archived       bool
	}{
		ProjectView:    pv,
		Eid:            eid,
		Ep:             ep,
		Tab:            tab,
		PageReportView: pageReportView,
		Archived:       isArchived && isTextMedia,
	}

	pageView := &PageView{
		Data:      data,
		User:      *user,
		PageTitle: "RESOURCES_VIEW_" + strings.ToUpper(tab),
	}

	h.Renderer.RenderTemplate(w, "resources", pageView)
}

// archiveHandler the HTTP request for the archive page. It loads the data from the
// archive and displays the source code of the crawler's response for a specific resource.
func (h *resourceHandler) archiveHandler(w http.ResponseWriter, r *http.Request) {
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

	rid, err := strconv.Atoi(r.URL.Query().Get("rid"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	eid := r.URL.Query().Get("eid")
	ep := r.URL.Query().Get("ep")
	if eid == "" && ep == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tab := r.URL.Query().Get("t")
	if tab == "" {
		tab = "details"
	}

	pv, err := h.ProjectViewService.GetProjectView(pid, user.Id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	isArchived := h.Container.ArchiveService.ArchiveExists(&pv.Project)
	if !isArchived {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pageReportView := h.ReportService.GetPageReport(rid, pv.Crawl.Id, "default", 1)
	isTextMedia := strings.HasPrefix(pageReportView.PageReport.MediaType, "text/")
	if !isTextMedia {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	record := h.Container.ArchiveService.ReadArchiveRecord(&pv.Project, pageReportView.PageReport.URL)

	data := &struct {
		PageReportView *models.PageReportView
		ProjectView    *models.ProjectView
		Eid            string
		Ep             string
		Tab            string
		ArchiveRecord  *models.ArchiveRecord
	}{
		ProjectView:    pv,
		PageReportView: pageReportView,
		Eid:            eid,
		Ep:             ep,
		Tab:            tab,
		ArchiveRecord:  record,
	}

	pageView := &PageView{
		Data:      data,
		User:      *user,
		PageTitle: "ARCHIVE_VIEW",
	}

	h.Renderer.RenderTemplate(w, "archive", pageView)
}
