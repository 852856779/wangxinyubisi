package page

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"strconv"
)

const tmpl string = `
{{ if .HasPages }}
<nav>
<ul class="pagination pagination-dark pagination-sm">
  {{ .GetPreviousButton "上一页" }}
  {{ range .FirstPart }}
    {{ . }}
  {{ end }}
  {{ if len .MiddlePart }}
    {{ .GetDots }}
	{{ range .MiddlePart }}
	  {{ . }}
	{{ end }}
  {{ end }}
  {{ if len .LastPart }}
	{{ .GetDots }}
	{{ range .LastPart }}
	  {{ . }}
	{{ end }}
  {{ end }}
  {{ .GetNextButton "下一页" }}
</ul>
</nav>
{{end}}
`

const onEachSide int = 3

func (p *Pagination) FirstPart() []string {
	return p.firstPart
}
func (p *Pagination) MiddlePart() []string {
	return p.middlePart
}
func (p *Pagination) LastPart() []string {
	return p.lastPart
}
func (p *Pagination) generate() {
	if !p.HasPages() {
		return
	}
	if p.TotalPages() < (onEachSide*2 + 6) { // 11页以内
		p.firstPart = p.getUrlRange(1, p.TotalPages())
	} else {
		window := onEachSide * 2
		lastPage := p.TotalPages()
		if p.currentPage < window { // 靠近开头
			p.firstPart = p.getUrlRange(1, window+2)
			p.lastPart = p.getUrlRange(lastPage-1, lastPage)
		} else if p.currentPage > (lastPage - window) { // 靠近结尾
			p.firstPart = p.getUrlRange(1, 2)
			p.lastPart = p.getUrlRange(lastPage-(window+2), lastPage)
		} else { // 在中间
			p.firstPart = p.getUrlRange(1, 2)
			p.middlePart = p.getUrlRange(p.currentPage-onEachSide, p.currentPage+onEachSide)
			p.lastPart = p.getUrlRange(lastPage-1, lastPage)
		}
	}
}

func (p *Pagination) getUrlRange(start, end int) []string {
	var ret []string
	for i := start; i <= end; i++ {
		ret = append(ret, p.getUrl(i, strconv.Itoa(i)))
	}
	return ret
}

func (p *Pagination) getUrl(page int, text string) string {
	strPage := strconv.Itoa(page)
	if p.currentPage == page {
		return p.GetActivePageWrapper(strPage)
	} else {
		baseUrl, _ := url.Parse(p.baseUrl)
		params := baseUrl.Query()
		delete(params, "page")
		strParam := ""
		for k, v := range params {
			strParam = strParam + "&" + k + "=" + v[0] // TODO
		}

		href := baseUrl.Host + "?page=" + strPage + strParam
		return p.GetAvailablePageWrapper(href, text)
	}
}

func (p *Pagination) GetActivePageWrapper(text string) string {
	return "<li class=\"page-item bg-dark\"><a class=\"page-link bg-dark text-white\">" + text + "</a></li>"
}
func (p *Pagination) GetDisabledPageWrapper(text string) string {
	return "<li class=\"page-item disabled\"><a class=\"page-link text-secondary\">" + text + "</a></li>"
}
func (p *Pagination) GetAvailablePageWrapper(href, page string) string {
	return "<li class=\"page-item\"><a class=\"page-link text-dark\" href=\"" + href + "\">" + page + "</a></li>"
}
func (p *Pagination) GetDots() string {
	return "<li class=\"page-item disabled\"><a class=\"page-link text-secondary\">...</a></li>"
}
func (p *Pagination) GetPreviousButton(text string) string { // "&laquo;"
	if p.currentPage <= 1 {
		return p.GetDisabledPageWrapper(text)
	}

	return p.getUrl(p.currentPage-1, "上一页")
}
func (p *Pagination) GetNextButton(text string) string { // &raquo;
	if p.currentPage == p.TotalPages() {
		return p.GetDisabledPageWrapper(text)
	}
	return p.getUrl(p.currentPage+1, "下一页")
}

// Render 生成html
func (p *Pagination) Render() template.HTML {
	p.generate()

	var out bytes.Buffer
	t := template.Must(template.New("pagination").Parse(tmpl))
	err := t.Execute(&out, p)
	if err != nil {
		return template.HTML(fmt.Sprintf("Error executing pagination template: %s", err))
	}
	return template.HTML(html.UnescapeString(out.String()))
}
