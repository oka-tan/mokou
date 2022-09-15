package asagi

import (
	"mokou/config"
	"mokou/utils"
	"regexp"
	"strings"
)

//postRestorer restores the API html
//from an Asagi post.
type postRestorer struct {
	quoteLink      *regexp.Regexp
	crossBoardLink *regexp.Regexp
	greenText      *regexp.Regexp
	codeTags       *regexp.Regexp
	bannedTags     *regexp.Regexp
	spoilerTags    *regexp.Regexp
	fortuneTags    *regexp.Regexp
	literalTags    *regexp.Regexp
	mootTags       *regexp.Regexp
	boldTags       *regexp.Regexp
	shiftJisTags   *regexp.Regexp

	enableCode    bool
	enableSpoiler bool
	enableFortune bool
	enableExif    bool
	enableOekaki  bool
}

//restoreComment restores a comment's HTML
//from the Asagi post.
func (p *postRestorer) restoreComment(post *post, exif map[string]interface{}) *string {
	if post.Comment == nil || *post.Comment == "" {
		return nil
	}

	rawComment, exifHasRawComment := exif["comment"]
	if exifHasRawComment {
		rawCommentString := rawComment.(string)
		return &rawCommentString
	}

	com := *post.Comment
	com = utils.FilterHtml(com)
	com = p.quoteLink.ReplaceAllString(com, "<a class=\"quotelink\" href=\"#p$1\">&gt;&gt;$1</a>")
	com = p.crossBoardLink.ReplaceAllString(com, "<a class=\"quotelink\" href=\"/$1/post/$2\">&gt;&gt;&gt;/$1/$2</a>")
	com = p.greenText.ReplaceAllString(com, "<span class=\"quote\">&gt;$1</span>")
	com = p.bannedTags.ReplaceAllString(com, "<strong style=\"color:red\">$1</strong>")
	com = p.mootTags.ReplaceAllString(com, "<div style=\"padding: 5px;margin-left: .5em;border-color: #faa;bprder: 2px dashed rgba(255,0,0,.1);border-radius: 2px\">$1</div>")
	com = p.boldTags.ReplaceAllString(com, "<strong>$1</strong>")
	com = p.shiftJisTags.ReplaceAllString(com, "<span class=\"sjis\">$1</span>")

	if p.enableCode {
		com = p.codeTags.ReplaceAllString(com, "<pre>$1</pre>")
	}
	if p.enableSpoiler {
		com = p.spoilerTags.ReplaceAllString(com, "<s>$1</s>")
	}
	if p.enableFortune {
		com = p.fortuneTags.ReplaceAllString(com, "<span class=\"fortune\" style=\"color:$1\"><br><br><b>$2</b></span>")
	}

	com = p.literalTags.ReplaceAllString(com, "[$1]")
	com = strings.ReplaceAll(com, "\n", "<br>")

	if p.enableExif {
		if hasExif(exif) {
			com += createExifTable(exif, post.Num)
		}
	}

	if p.enableOekaki {
		if hasOekaki(exif) {
			com += createOekaki(exif)
		}
	}

	return &com
}

//newPostRestorer constructs a PostRestorer from a board configuration
func newPostRestorer(boardConfig *config.AsagiBoardConfig) postRestorer {
	//All of the regex is written under the assumptions that:
	//", &, < and > have been converted to &quot;
	//&amp; &lt; and &gt;
	//But \n remains as \n
	//Converting \n last is convenient for the greentext regex.
	quoteLink := regexp.MustCompile("&gt;&gt;(\\d+)")
	crossBoardLink := regexp.MustCompile("&gt;&gt;&gt;\\/(\\w+)\\/(\\w+)?")
	greenText := regexp.MustCompile("(?m)^&gt;(.*)$")
	bannedTags := regexp.MustCompile("\\[banned\\](.*)\\[/banned\\]")
	codeTags := regexp.MustCompile("\\[code\\](.*?)\\[/code\\]")
	spoilerTags := regexp.MustCompile("\\[spoiler\\](.*?)\\[/spoiler\\]")
	fortuneTags := regexp.MustCompile("\n\n\\[fortune color=&quot;(.*?)&quot;\\](.*)\\[\\/fortune\\]")
	literalTags := regexp.MustCompile("\\[(\\S*?):lit\\]")
	mootTags := regexp.MustCompile("\\[moot\\](.*?)\\[moot\\]")
	boldTags := regexp.MustCompile("\\[b\\](.*?)\\[\\/b\\]")
	shiftJisTags := regexp.MustCompile("\\[shiftjis\\](.*?)\\[\\/shiftjis\\]")

	return postRestorer{
		quoteLink:      quoteLink,
		crossBoardLink: crossBoardLink,
		greenText:      greenText,
		codeTags:       codeTags,
		bannedTags:     bannedTags,
		spoilerTags:    spoilerTags,
		fortuneTags:    fortuneTags,
		literalTags:    literalTags,
		mootTags:       mootTags,
		boldTags:       boldTags,
		shiftJisTags:   shiftJisTags,

		enableCode:    boardConfig.EnableCode,
		enableFortune: boardConfig.EnableFortune,
		enableSpoiler: boardConfig.EnableSpoiler,
		enableExif:    boardConfig.EnableExif,
		enableOekaki:  boardConfig.EnableOekaki,
	}
}
