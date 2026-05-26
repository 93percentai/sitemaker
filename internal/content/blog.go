package content

import (
	"sort"
)

func (s *Site) ClassifyPages() {
	for _, p := range s.Pages {
		if p.IsBlogPost && !p.Draft {
			s.Posts = append(s.Posts, p)
		}
	}

	sort.Slice(s.Posts, func(i, j int) bool {
		return s.Posts[i].Date.After(s.Posts[j].Date)
	})
}

func (s *Site) BuildTagIndex() {
	s.TagIndex = make(map[string][]*Page)
	for _, p := range s.Posts {
		for _, tag := range p.Tags {
			s.TagIndex[tag] = append(s.TagIndex[tag], p)
		}
	}
}

func (s *Site) LinkPrevNext() {
	for i, p := range s.Posts {
		if i > 0 {
			p.Next = s.Posts[i-1]
		}
		if i < len(s.Posts)-1 {
			p.Prev = s.Posts[i+1]
		}
	}
}
