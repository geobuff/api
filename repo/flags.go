package repo

type FlagEntry struct {
	Code string `json:"code"`
	Url  string `json:"url"`
}

type FlagGroup struct {
	Key     string      `json:"key"`
	Label   string      `json:"label"`
	Entries []FlagEntry `json:"entries"`
}

var FlagGroups = []FlagGroup{
	{
		Key:   "world",
		Label: "🌎 World, Countries",
		Entries: []FlagEntry{
			{
				Code: "buff",
				Url:  "https://ik.imagekit.io/ucszu5sud3vz/flag-geobuff-sm",
			},
			{
				Code: "ad",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1e9.svg",
			},
			{
				Code: "ae",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1ea.svg",
			},
			{
				Code: "af",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1eb.svg",
			},
			{
				Code: "ag",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1ec.svg",
			},
			{
				Code: "al",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f1.svg",
			},
			{
				Code: "am",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f2.svg",
			},
			{
				Code: "ao",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f4.svg",
			},
			{
				Code: "ar",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg",
			},
			{
				Code: "at",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f9.svg",
			},
			{
				Code: "au",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg",
			},
			{
				Code: "az",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1ff.svg",
			},
			{
				Code: "ba",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1e6.svg",
			},
			{
				Code: "bb",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1e7.svg",
			},
			{
				Code: "bd",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1e9.svg",
			},
			{
				Code: "be",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ea.svg",
			},
			{
				Code: "bf",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1eb.svg",
			},
			{
				Code: "bg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ec.svg",
			},
			{
				Code: "bh",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ed.svg",
			},
			{
				Code: "bi",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ee.svg",
			},
			{
				Code: "bj",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ef.svg",
			},
			{
				Code: "bn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f3.svg",
			},
			{
				Code: "bo",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f4.svg",
			},
			{
				Code: "br",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg",
			},
			{
				Code: "bs",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f8.svg",
			},
			{
				Code: "bt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f9.svg",
			},
			{
				Code: "bw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1fc.svg",
			},
			{
				Code: "by",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1fe.svg",
			},
			{
				Code: "bz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1ff.svg",
			},
			{
				Code: "ca",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg",
			},
			{
				Code: "cd",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e9.svg",
			},
			{
				Code: "cf",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1eb.svg",
			},
			{
				Code: "cg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1ec.svg",
			},
			{
				Code: "ch",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1ed.svg",
			},
			{
				Code: "ci",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1ee.svg",
			},
			{
				Code: "cl",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f1.svg",
			},
			{
				Code: "cm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f2.svg",
			},
			{
				Code: "cn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f3.svg",
			},
			{
				Code: "co",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f4.svg",
			},
			{
				Code: "cr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f7.svg",
			},
			{
				Code: "cu",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1fa.svg",
			},
			{
				Code: "cv",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1fb.svg",
			},
			{
				Code: "cy",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1fe.svg",
			},
			{
				Code: "cz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1ff.svg",
			},
			{
				Code: "de",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg",
			},
			{
				Code: "dj",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ef.svg",
			},
			{
				Code: "dk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1f0.svg",
			},
			{
				Code: "dm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1f2.svg",
			},
			{
				Code: "do",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1f4.svg",
			},
			{
				Code: "dz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ff.svg",
			},
			{
				Code: "ec",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1e8.svg",
			},
			{
				Code: "ee",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1ea.svg",
			},
			{
				Code: "eg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1ec.svg",
			},
			{
				Code: "er",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f7.svg",
			},
			{
				Code: "es",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f8.svg",
			},
			{
				Code: "et",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f9.svg",
			},
			{
				Code: "fi",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1ee.svg",
			},
			{
				Code: "fj",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1ef.svg",
			},
			{
				Code: "fm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f2.svg",
			},
			{
				Code: "fr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg",
			},
			{
				Code: "ga",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e6.svg",
			},
			{
				Code: "gb",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg",
			},
			{
				Code: "gd",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e9.svg",
			},
			{
				Code: "ge",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1ea.svg",
			},
			{
				Code: "gh",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1ed.svg",
			},
			{
				Code: "gm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1f2.svg",
			},
			{
				Code: "gn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1f3.svg",
			},
			{
				Code: "gq",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1f6.svg",
			},
			{
				Code: "gr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1f7.svg",
			},
			{
				Code: "gt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1f9.svg",
			},
			{
				Code: "gw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1fc.svg",
			},
			{
				Code: "gy",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1fe.svg",
			},
			{
				Code: "hn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ed-1f1f3.svg",
			},
			{
				Code: "hr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ed-1f1f7.svg",
			},
			{
				Code: "ht",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ed-1f1f9.svg",
			},
			{
				Code: "hu",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ed-1f1fa.svg",
			},
			{
				Code: "id",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1e9.svg",
			},
			{
				Code: "ie",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1ea.svg",
			},
			{
				Code: "il",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f1.svg",
			},
			{
				Code: "in",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f3.svg",
			},
			{
				Code: "iq",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f6.svg",
			},
			{
				Code: "ir",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f7.svg",
			},
			{
				Code: "is",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f8.svg",
			},
			{
				Code: "it",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f9.svg",
			},
			{
				Code: "jm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ef-1f1f2.svg",
			},
			{
				Code: "jo",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ef-1f1f4.svg",
			},
			{
				Code: "jp",
				Url:  "https://ik.imagekit.io/ucszu5sud3vz/flag-jp_9E4mwVsqHDW.svg",
			},
			{
				Code: "ke",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1ea.svg",
			},
			{
				Code: "kg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1ec.svg",
			},
			{
				Code: "kh",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1ed.svg",
			},
			{
				Code: "ki",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1ee.svg",
			},
			{
				Code: "km",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1f2.svg",
			},
			{
				Code: "kn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1f3.svg",
			},
			{
				Code: "kp",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1f5.svg",
			},
			{
				Code: "kr",
				Url:  "https://ik.imagekit.io/ucszu5sud3vz/flag-kr_tYV56l4dBQ3.svg",
			},
			{
				Code: "kw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1fc.svg",
			},
			{
				Code: "kz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f0-1f1ff.svg",
			},
			{
				Code: "la",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1e6.svg",
			},
			{
				Code: "lb",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1e7.svg",
			},
			{
				Code: "lc",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1e8.svg",
			},
			{
				Code: "li",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1ee.svg",
			},
			{
				Code: "lk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1f0.svg",
			},
			{
				Code: "lr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1f7.svg",
			},
			{
				Code: "ls",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1f8.svg",
			},
			{
				Code: "lt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1f9.svg",
			},
			{
				Code: "lu",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1fa.svg",
			},
			{
				Code: "lv",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1fb.svg",
			},
			{
				Code: "ly",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f1-1f1fe.svg",
			},
			{
				Code: "ma",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1e6.svg",
			},
			{
				Code: "mc",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1e8.svg",
			},
			{
				Code: "md",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1e9.svg",
			},
			{
				Code: "me",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1ea.svg",
			},
			{
				Code: "mg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1ec.svg",
			},
			{
				Code: "mh",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1ed.svg",
			},
			{
				Code: "mk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f0.svg",
			},
			{
				Code: "ml",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f1.svg",
			},
			{
				Code: "mm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f2.svg",
			},
			{
				Code: "mn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f3.svg",
			},
			{
				Code: "mr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f7.svg",
			},
			{
				Code: "mt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1f9.svg",
			},
			{
				Code: "mu",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fa.svg",
			},
			{
				Code: "mv",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fb.svg",
			},
			{
				Code: "mw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fc.svg",
			},
			{
				Code: "mx",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fd.svg",
			},
			{
				Code: "my",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fe.svg",
			},
			{
				Code: "mz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1ff.svg",
			},
			{
				Code: "na",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1e6.svg",
			},
			{
				Code: "ne",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ea.svg",
			},
			{
				Code: "ng",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ec.svg",
			},
			{
				Code: "ni",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ee.svg",
			},
			{
				Code: "nl",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1f1.svg",
			},
			{
				Code: "no",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1f4.svg",
			},
			{
				Code: "np",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1f5.svg",
			},
			{
				Code: "nr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1f7.svg",
			},
			{
				Code: "nz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ff.svg",
			},
			{
				Code: "om",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f4-1f1f2.svg",
			},
			{
				Code: "pa",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1e6.svg",
			},
			{
				Code: "pe",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1ea.svg",
			},
			{
				Code: "pg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1ec.svg",
			},
			{
				Code: "ph",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1ed.svg",
			},
			{
				Code: "pk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1f0.svg",
			},
			{
				Code: "pl",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1f1.svg",
			},
			{
				Code: "ps",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1f8.svg",
			},
			{
				Code: "pt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1f9.svg",
			},
			{
				Code: "pw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1fc.svg",
			},
			{
				Code: "py",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1fe.svg",
			},
			{
				Code: "qa",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f6-1f1e6.svg",
			},
			{
				Code: "ro",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f7-1f1f4.svg",
			},
			{
				Code: "rs",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f7-1f1f8.svg",
			},
			{
				Code: "ru",
				Url:  "https://ik.imagekit.io/ucszu5sud3vz/flag-ru_hI7yCCaMi7E.svg",
			},
			{
				Code: "rw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f7-1f1fc.svg",
			},
			{
				Code: "sa",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1e6.svg",
			},
			{
				Code: "sb",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1e7.svg",
			},
			{
				Code: "sc",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1e8.svg",
			},
			{
				Code: "sd",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1e9.svg",
			},
			{
				Code: "se",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1ea.svg",
			},
			{
				Code: "sg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1ec.svg",
			},
			{
				Code: "si",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1ee.svg",
			},
			{
				Code: "sk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f0.svg",
			},
			{
				Code: "sl",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f1.svg",
			},
			{
				Code: "sm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f2.svg",
			},
			{
				Code: "sn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f3.svg",
			},
			{
				Code: "so",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f4.svg",
			},
			{
				Code: "sr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f7.svg",
			},
			{
				Code: "ss",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f8.svg",
			},
			{
				Code: "st",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1f9.svg",
			},
			{
				Code: "sv",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1fb.svg",
			},
			{
				Code: "sy",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1fe.svg",
			},
			{
				Code: "sz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f8-1f1ff.svg",
			},
			{
				Code: "td",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1e9.svg",
			},
			{
				Code: "tg",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1ec.svg",
			},
			{
				Code: "th",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1ed.svg",
			},
			{
				Code: "tj",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1ef.svg",
			},
			{
				Code: "tl",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f1.svg",
			},
			{
				Code: "tm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f2.svg",
			},
			{
				Code: "tn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f3.svg",
			},
			{
				Code: "to",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f4.svg",
			},
			{
				Code: "tr",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f7.svg",
			},
			{
				Code: "tt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f9.svg",
			},
			{
				Code: "tv",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1fb.svg",
			},
			{
				Code: "tw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1fc.svg",
			},
			{
				Code: "tz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1ff.svg",
			},
			{
				Code: "ua",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1e6.svg",
			},
			{
				Code: "ug",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1ec.svg",
			},
			{
				Code: "us",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg",
			},
			{
				Code: "uy",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1fe.svg",
			},
			{
				Code: "uz",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1ff.svg",
			},
			{
				Code: "va",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fb-1f1e6.svg",
			},
			{
				Code: "vc",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fb-1f1e8.svg",
			},
			{
				Code: "ve",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fb-1f1ea.svg",
			},
			{
				Code: "vn",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fb-1f1f3.svg",
			},
			{
				Code: "vu",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fb-1f1fa.svg",
			},
			{
				Code: "ws",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fc-1f1f8.svg",
			},
			{
				Code: "xk",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fd-1f1f0.svg",
			},
			{
				Code: "ye",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fe-1f1ea.svg",
			},
			{
				Code: "za",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1e6.svg",
			},
			{
				Code: "zm",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1f2.svg",
			},
			{
				Code: "zw",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1fc.svg",
			},
		},
	},
	{
		Key:   "au",
		Label: "🇦🇺 Australia, States and Territories",
		Entries: []FlagEntry{
			{
				Code: "au-nsw",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Flag_of_New_South_Wales.svg",
			},
			{
				Code: "au-qld",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/04/Flag_of_Queensland.svg",
			},
			{
				Code: "au-sa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fd/Flag_of_South_Australia.svg",
			},
			{
				Code: "au-tas",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/46/Flag_of_Tasmania.svg",
			},
			{
				Code: "au-vic",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/08/Flag_of_Victoria_%28Australia%29.svg",
			},
			{
				Code: "au-wa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a5/Flag_of_Western_Australia.svg",
			},
			{
				Code: "au-act",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8c/Flag_of_the_Australian_Capital_Territory.svg",
			},
			{
				Code: "au-jbt",
				Url:  "https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg",
			},
			{
				Code: "au-nt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b7/Flag_of_the_Northern_Territory.svg",
			},
		},
	},
	{
		Key:   "ar",
		Label: "🇦🇷 Argentina, Provinces",
		Entries: []FlagEntry{
			{
				Code: "ar-ba",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/15/Bandera_de_la_Provincia_de_Buenos_Aires.svg",
			},
			{
				Code: "ar-ct",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7b/Bandera_de_la_Provincia_de_Catamarca.svg",
			},
			{
				Code: "ar-cc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/33/Bandera_de_la_Provincia_del_Chaco.svg",
			},
			{
				Code: "ar-ch",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/88/Bandera_de_la_Provincia_del_Chubut.svg",
			},
			{
				Code: "ar-df",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f5/Bandera_de_la_Ciudad_de_Buenos_Aires.svg",
			},
			{
				Code: "ar-cb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/96/Bandera_de_la_Provincia_de_C%C3%B3rdoba.svg",
			},
			{
				Code: "ar-cn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/46/Bandera_de_la_Provincia_de_Corrientes.svg",
			},
			{
				Code: "ar-er",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5b/Bandera_de_la_Provincia_de_Entre_R%C3%ADos.svg",
			},
			{
				Code: "ar-fm",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/42/Bandera_de_la_Provincia_de_Formosa.svg",
			},
			{
				Code: "ar-jy",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c9/Bandera_de_la_Provincia_de_Jujuy.svg",
			},
			{
				Code: "ar-lp",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/81/Bandera_de_la_Provincia_de_La_Pampa.svg",
			},
			{
				Code: "ar-lr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/60/Bandera_de_la_Provincia_de_La_Rioja.svg",
			},
			{
				Code: "ar-mz",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7c/Bandera_de_la_Provincia_de_Mendoza.svg",
			},
			{
				Code: "ar-mn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ce/Bandera_de_la_Provincia_de_Misiones.svg",
			},
			{
				Code: "ar-nq",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bc/Bandera_de_la_Provincia_de_Neuqu%C3%A9n.svg",
			},
			{
				Code: "ar-rn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5d/Bandera_de_la_Provincia_del_R%C3%ADo_Negro.svg",
			},
			{
				Code: "ar-sa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6a/Bandera_de_la_Provincia_de_Salta.svg",
			},
			{
				Code: "ar-sj",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8a/Bandera_de_San_Juan_Ciudadana.png",
			},
			{
				Code: "ar-sl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0e/Bandera_de_la_Provincia_de_San_Luis.svg",
			},
			{
				Code: "ar-sc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/45/Bandera_de_la_Provincia_de_Santa_Cruz.svg",
			},
			{
				Code: "ar-sf",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/84/Bandera_de_la_Provincia_de_Santa_Fe.svg",
			},
			{
				Code: "ar-se",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/07/Bandera_de_la_Provincia_de_Santiago_del_Estero.svg",
			},
			{
				Code: "ar-tf",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/94/Bandera_de_la_Provincia_de_Tierra_del_Fuego.svg",
			},
			{
				Code: "ar-tm",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ce/Bandera_de_la_Provincia_de_Tucum%C3%A1n.svg",
			},
		},
	},
	{
		Key:   "br",
		Label: "🇧🇷 Brazil, States",
		Entries: []FlagEntry{
			{
				Code: "br-am",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6b/Bandeira_do_Amazonas.svg",
			},
			{
				Code: "br-pa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/02/Bandeira_do_Par%C3%A1.svg",
			},
			{
				Code: "br-mt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0b/Bandeira_de_Mato_Grosso.svg",
			},
			{
				Code: "br-mg",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f4/Bandeira_de_Minas_Gerais.svg",
			},
			{
				Code: "br-ba",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/28/Bandeira_da_Bahia.svg",
			},
			{
				Code: "br-ms",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/64/Bandeira_de_Mato_Grosso_do_Sul.svg",
			},
			{
				Code: "br-go",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/be/Flag_of_Goi%C3%A1s.svg",
			},
			{
				Code: "br-ma",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/45/Bandeira_do_Maranh%C3%A3o.svg",
			},
			{
				Code: "br-rs",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/63/Bandeira_do_Rio_Grande_do_Sul.svg",
			},
			{
				Code: "br-to",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/ff/Bandeira_do_Tocantins.svg",
			},
			{
				Code: "br-pi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/33/Bandeira_do_Piau%C3%AD.svg",
			},
			{
				Code: "br-sp",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2b/Bandeira_do_estado_de_S%C3%A3o_Paulo.svg",
			},
			{
				Code: "br-ro",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fa/Bandeira_de_Rond%C3%B4nia.svg",
			},
			{
				Code: "br-rr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/98/Bandeira_de_Roraima.svg",
			},
			{
				Code: "br-pr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/93/Bandeira_do_Paran%C3%A1.svg",
			},
			{
				Code: "br-ac",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4c/Bandeira_do_Acre.svg",
			},
			{
				Code: "br-ce",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2e/Bandeira_do_Cear%C3%A1.svg",
			},
			{
				Code: "br-ap",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0c/Bandeira_do_Amap%C3%A1.svg",
			},
			{
				Code: "br-pe",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/59/Bandeira_de_Pernambuco.svg",
			},
			{
				Code: "br-sc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1a/Bandeira_de_Santa_Catarina.svg",
			},
			{
				Code: "br-pb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bb/Bandeira_da_Para%C3%ADba.svg",
			},
			{
				Code: "br-rn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/30/Bandeira_do_Rio_Grande_do_Norte.svg",
			},
			{
				Code: "br-es",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/43/Bandeira_do_Esp%C3%ADrito_Santo.svg",
			},
			{
				Code: "br-rj",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/73/Bandeira_do_estado_do_Rio_de_Janeiro.svg",
			},
			{
				Code: "br-al",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/88/Bandeira_de_Alagoas.svg",
			},
			{
				Code: "br-se",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/be/Bandeira_de_Sergipe.svg",
			},
			{
				Code: "br-df",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3c/Bandeira_do_Distrito_Federal_%28Brasil%29.svg",
			},
		},
	},
	{
		Key:   "ca",
		Label: "🇨🇦 Canada, Provinces and Territories",
		Entries: []FlagEntry{

			{
				Code: "ca-nu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/90/Flag_of_Nunavut.svg",
			},
			{
				Code: "ca-qc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5f/Flag_of_Quebec.svg",
			},
			{
				Code: "ca-nt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c1/Flag_of_the_Northwest_Territories.svg",
			},
			{
				Code: "ca-on",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/88/Flag_of_Ontario.svg",
			},
			{
				Code: "ca-bc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_of_British_Columbia.svg",
			},
			{
				Code: "ca-ab",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f5/Flag_of_Alberta.svg",
			},
			{
				Code: "ca-sk",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bb/Flag_of_Saskatchewan.svg",
			},
			{
				Code: "ca-mb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c4/Flag_of_Manitoba.svg",
			},
			{
				Code: "ca-yt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/69/Flag_of_Yukon.svg",
			},
			{
				Code: "ca-nl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/dd/Flag_of_Newfoundland_and_Labrador.svg",
			},
			{
				Code: "ca-nb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fb/Flag_of_New_Brunswick.svg",
			},
			{
				Code: "ca-ns",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c0/Flag_of_Nova_Scotia.svg",
			},
			{
				Code: "ca-pe",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d7/Flag_of_Prince_Edward_Island.svg",
			},
		},
	},
	{
		Key:   "co",
		Label: "🇨🇴 Colombia, Departments",
		Entries: []FlagEntry{

			{
				Code: "co-ama",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1f/Flag_of_Amazonas_%28Colombia%29.svg",
			},
			{
				Code: "co-ant",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/33/Flag_of_Antioquia_Department.svg",
			},
			{
				Code: "co-ara",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8a/Flag_of_Arauca.svg",
			},
			{
				Code: "co-atl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a3/Flag_of_Atl%C3%A1ntico.svg",
			},
			{
				Code: "co-bol",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a0/Flag_of_Bol%C3%ADvar_%28Colombia%29.svg",
			},
			{
				Code: "co-boy",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2e/Flag_of_Boyac%C3%A1_Department.svg",
			},
			{
				Code: "co-cau",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ae/Flag_of_Cauca.svg",
			},
			{
				Code: "co-ces",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/43/Flag_of_Cesar.svg",
			},
			{
				Code: "co-cho",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ae/Flag_of_Choc%C3%B3.svg",
			},
			{
				Code: "co-cal",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e2/Flag_of_Caldas.svg",
			},
			{
				Code: "co-cor",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9b/Flag_of_C%C3%B3rdoba.svg",
			},
			{
				Code: "co-caq",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b0/Flag_of_Caquet%C3%A1.svg",
			},
			{
				Code: "co-cas",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/ba/Flag_of_Casanare.svg",
			},
			{
				Code: "co-cun",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Flag_of_Cundinamarca.svg",
			},
			{
				Code: "co-dc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9e/Flag_of_Bogot%C3%A1.svg",
			},
			{
				Code: "co-gua",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Flag_of_Guain%C3%ADa.svg",
			},
			{
				Code: "co-guv",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9b/Flag_of_Guaviare.svg",
			},
			{
				Code: "co-hui",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b2/Flag_of_Huila.svg",
			},
			{
				Code: "co-lag",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/35/Flag_of_La_Guajira.svg",
			},
			{
				Code: "co-mag",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/18/Flag_of_Magdalena.svg",
			},
			{
				Code: "co-met",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/72/Flag_of_Meta.svg",
			},
			{
				Code: "co-nar",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/20/Flag_of_Nari%C3%B1o.svg",
			},
			{
				Code: "co-nsa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9a/Flag_of_Norte_de_Santander.svg",
			},
			{
				Code: "co-put",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/93/Flag_of_Putumayo.svg",
			},
			{
				Code: "co-qui",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/29/Flag_of_Quind%C3%ADo.svg",
			},
			{
				Code: "co-ris",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/02/Flag_of_Risaralda.svg",
			},
			{
				Code: "co-san",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6f/Flag_of_Santander_%28Colombia%29.svg",
			},
			{
				Code: "co-suc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/39/Flag_of_Sucre_%28Colombia%29.svg",
			},
			{
				Code: "co-sap",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/36/Flag_of_San_Andr%C3%A9s_y_Providencia.svg",
			},
			{
				Code: "co-tol",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b9/Flag_of_Tolima.svg",
			},
			{
				Code: "co-vac",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d3/Flag_of_Valle_del_Cauca.svg",
			},
			{
				Code: "co-vid",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fc/Flag_of_Vichada.svg",
			},
			{
				Code: "co-vau",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c0/Flag_of_Vaup%C3%A9s.svg",
			},
		},
	},
	{
		Key:   "fr",
		Label: "🇫🇷 France, Regions",
		Entries: []FlagEntry{

			{
				Code: "fr-ar",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/85/Flag_of_the_region_Auvergne-Rh%C3%B4ne-Alpes.svg",
			},
			{
				Code: "fr-bf",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/28/Flag_of_the_region_Bourgogne-Franche-Comt%C3%A9.svg",
			},
			{
				Code: "fr-bt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/29/Flag_of_Brittany_%28Gwenn_ha_du%29.svg",
			},
			{
				Code: "fr-cn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d5/Flag_of_Centre-Val_de_Loire.svg",
			},
			{
				Code: "fr-ce",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7c/Flag_of_Corsica.svg",
			},
			{
				Code: "fr-ao",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f3/Proposed_design_for_a_flag_of_Grand_Est.svg",
			},
			{
				Code: "fr-nc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bd/Proposed_design_for_a_flag_of_Hauts-de-France.svg",
			},
			{
				Code: "fr-if",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3b/Proposed_flag_of_%C3%8Ele-de-France.svg",
			},
			{
				Code: "fr-nd",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/41/Flag_of_Normandie.svg",
			},
			{
				Code: "fr-ac",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b6/Flag_of_Nouvelle-Aquitaine.svg",
			},
			{
				Code: "fr-lp",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1d/Flag_of_R%C3%A9gion_Occitanie_%28symbol_only%29.svg",
			},
			{
				Code: "fr-pl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Unofficial_flag_of_Pays-de-la-Loire.svg",
			},
			{
				Code: "fr-pr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/94/Flag_of_Provence-Alpes-C%C3%B4te_d%27Azur.svg",
			},
		},
	},
	{
		Key:   "de",
		Label: "🇩🇪 Germany, States",
		Entries: []FlagEntry{

			{
				Code: "de-bw",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5c/Flag_of_Baden-W%C3%BCrttemberg.svg",
			},
			{
				Code: "de-by",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/16/Flag_of_Bavaria_%28striped%29.svg",
			},
			{
				Code: "de-be",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ec/Flag_of_Berlin.svg",
			},
			{
				Code: "de-bb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/01/Flag_of_Brandenburg.svg",
			},
			{
				Code: "de-hb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0e/Flag_of_Bremen.svg",
			},
			{
				Code: "de-hh",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/74/Flag_of_Hamburg.svg",
			},
			{
				Code: "de-he",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Hesse.svg",
			},
			{
				Code: "de-ni",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/81/Flag_of_Lower_Saxony.svg",
			},
			{
				Code: "de-mv",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ce/Flag_of_Mecklenburg-Western_Pomerania.svg",
			},
			{
				Code: "de-nw",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/84/Flag_of_North_Rhine-Westphalia.svg",
			},
			{
				Code: "de-rp",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b7/Flag_of_Rhineland-Palatinate.svg",
			},
			{
				Code: "de-sl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Saarland.svg",
			},
			{
				Code: "de-sn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fd/Flag_of_Saxony.svg",
			},
			{
				Code: "de-st",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c2/Flag_of_Saxony-Anhalt_%28state%29.svg",
			},
			{
				Code: "de-sh",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b4/Flag_of_Schleswig-Holstein.svg",
			},
			{
				Code: "de-th",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bd/Flag_of_Thuringia.svg",
			},
		},
	},
	{
		Key:   "gb",
		Label: "🇬🇧 Great Britain, Countries",
		Entries: []FlagEntry{

			{
				Code: "gb-eng",
				Url:  "https://upload.wikimedia.org/wikipedia/en/b/be/Flag_of_England.svg",
			},
			{
				Code: "gb-sct",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/10/Flag_of_Scotland.svg",
			},
			{
				Code: "gb-wls",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/dc/Flag_of_Wales.svg",
			},
		},
	},
	{
		Key:   "it",
		Label: "🇮🇹 Italy, Regions",
		Entries: []FlagEntry{

			{
				Code: "it-65",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/45/Flag_of_Abruzzo.svg",
			},
			{
				Code: "it-77",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8e/Flag_of_Basilicata.svg",
			},
			{
				Code: "it-78",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8b/Flag_of_Calabria.svg",
			},
			{
				Code: "it-72",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c5/Flag_of_Campania.svg",
			},
			{
				Code: "it-45",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/77/Flag_of_Emilia-Romagna_%28de_facto%29.svg",
			},
			{
				Code: "it-36",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/55/Flag_of_Friuli-Venezia_Giulia.svg",
			},
			{
				Code: "it-62",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e1/Flag_of_Lazio.svg",
			},
			{
				Code: "it-42",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/88/Flag_of_Liguria.svg",
			},
			{
				Code: "it-25",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ea/Flag_of_Lombardy.svg",
			},
			{
				Code: "it-57",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/07/Flag_of_Marche.svg",
			},
			{
				Code: "it-67",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/84/Flag_of_Molise.svg",
			},
			{
				Code: "it-21",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b9/Flag_of_Piedmont.svg",
			},
			{
				Code: "it-75",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_of_Apulia.svg",
			},
			{
				Code: "it-88",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4e/Flag_of_Sardinia%2C_Italy.svg",
			},
			{
				Code: "it-82",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/84/Sicilian_Flag.svg",
			},
			{
				Code: "it-52",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f2/Flag_of_Tuscany.svg",
			},
			{
				Code: "it-32",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5f/Flag_of_South_Tyrol.svg",
			},
			{
				Code: "it-55",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/cc/Flag_of_Umbria.svg",
			},
			{
				Code: "it-23",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/90/Flag_of_Valle_d%27Aosta.svg",
			},
			{
				Code: "it-34",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d5/Flag_of_Veneto.svg",
			},
		},
	},
	{
		Key:   "jp",
		Label: "🇯🇵 Japan, Prefectures",
		Entries: []FlagEntry{

			{
				Code: "jp-01",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/22/Flag_of_Hokkaido_Prefecture.svg",
			},
			{
				Code: "jp-02",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/30/Flag_of_Aomori_Prefecture.svg",
			},
			{
				Code: "jp-03",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a9/Flag_of_Iwate_Prefecture.svg",
			},
			{
				Code: "jp-04",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c7/Flag_of_Miyagi_Prefecture.svg",
			},
			{
				Code: "jp-05",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/84/Flag_of_Akita_Prefecture.svg",
			},
			{
				Code: "jp-06",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a1/Flag_of_Yamagata_Prefecture.svg",
			},
			{
				Code: "jp-07",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4b/Flag_of_Fukushima_Prefecture.svg",
			},
			{
				Code: "jp-08",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a8/Flag_of_Ibaraki_Prefecture.svg",
			},
			{
				Code: "jp-09",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d5/Flag_of_Tochigi_Prefecture.svg",
			},
			{
				Code: "jp-10",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/ba/Flag_of_Gunma_Prefecture.svg",
			},
			{
				Code: "jp-11",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/cd/Flag_of_Saitama_Prefecture.svg",
			},
			{
				Code: "jp-12",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0a/Flag_of_Chiba_Prefecture.svg",
			},
			{
				Code: "jp-13",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/15/Flag_of_Tokyo_Metropolis.svg",
			},
			{
				Code: "jp-14",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a7/Flag_of_Kanagawa_Prefecture.svg",
			},
			{
				Code: "jp-15",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/cb/Flag_of_Niigata_Prefecture.svg",
			},
			{
				Code: "jp-16",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1d/Flag_of_Toyama_Prefecture.svg",
			},
			{
				Code: "jp-17",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6a/Flag_of_Ishikawa_Prefecture.svg",
			},
			{
				Code: "jp-18",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/56/Flag_of_Fukui_Prefecture.svg",
			},
			{
				Code: "jp-19",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Flag_of_Yamanashi_Prefecture.svg",
			},
			{
				Code: "jp-20",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f0/Flag_of_Nagano_Prefecture.svg",
			},
			{
				Code: "jp-21",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3e/Flag_of_Gifu_Prefecture.svg",
			},
			{
				Code: "jp-22",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/92/Flag_of_Shizuoka_Prefecture.svg",
			},
			{
				Code: "jp-23",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/02/Flag_of_Aichi_Prefecture.svg",
			},
			{
				Code: "jp-24",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8c/Flag_of_Mie_Prefecture.svg",
			},
			{
				Code: "jp-25",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/99/Flag_of_Shiga_Prefecture.svg",
			},
			{
				Code: "jp-26",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/06/Flag_of_Kyoto_Prefecture.svg",
			},
			{
				Code: "jp-27",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5a/Flag_of_Osaka_Prefecture.svg",
			},
			{
				Code: "jp-28",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/74/Flag_of_Hyogo_Prefecture.svg",
			},
			{
				Code: "jp-29",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Flag_of_Nara_Prefecture.svg",
			},
			{
				Code: "jp-30",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6e/Flag_of_Wakayama_Prefecture.svg",
			},
			{
				Code: "jp-31",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1c/Flag_of_Tottori_Prefecture.svg",
			},
			{
				Code: "jp-32",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e8/Flag_of_Shimane_Prefecture.svg",
			},
			{
				Code: "jp-33",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/33/Flag_of_Okayama_Prefecture.svg",
			},
			{
				Code: "jp-34",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ed/Flag_of_Hiroshima_Prefecture.svg",
			},
			{
				Code: "jp-35",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b9/Flag_of_Yamaguchi_Prefecture.svg",
			},
			{
				Code: "jp-36",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ac/Flag_of_Tokushima_Prefecture.svg",
			},
			{
				Code: "jp-37",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/29/Flag_of_Kagawa_Prefecture.svg",
			},
			{
				Code: "jp-38",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2d/Flag_of_Ehime_Prefecture.svg",
			},
			{
				Code: "jp-39",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/50/Flag_of_Kochi_Prefecture.svg",
			},
			{
				Code: "jp-40",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/71/Flag_of_Fukuoka_Prefecture.svg",
			},
			{
				Code: "jp-41",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/18/Flag_of_Saga_Prefecture.svg",
			},
			{
				Code: "jp-42",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Flag_of_Nagasaki_Prefecture.svg",
			},
			{
				Code: "jp-43",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Kumamoto_Prefecture.svg",
			},
			{
				Code: "jp-44",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c8/Flag_of_Oita_Prefecture.svg",
			},
			{
				Code: "jp-45",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0b/Flag_of_Miyazaki_Prefecture.svg",
			},
			{
				Code: "jp-46",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c5/Flag_of_Kagoshima_Prefecture.svg",
			},
			{
				Code: "jp-47",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2f/Flag_of_Okinawa_Prefecture.svg",
			},
		},
	},
	{
		Key:   "ru",
		Label: "🇷🇺 Russia, Federal Subjects",
		Entries: []FlagEntry{

			{
				Code: "ru-ad",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/16/Flag_of_Adygea.svg",
			},
			{
				Code: "ru-alt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/ff/Flag_of_Altai_Republic.svg",
			},
			{
				Code: "ru-amu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7a/Flag_of_Amur_Oblast.svg",
			},
			{
				Code: "ru-ark",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/16/Flag_of_Arkhangelsk_Oblast.svg",
			},
			{
				Code: "ru-ast",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/87/Flag_of_Astrakhan_Oblast.svg",
			},
			{
				Code: "ru-ba",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3d/Flag_of_Bashkortostan.svg",
			},
			{
				Code: "ru-bel",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2b/Flag_of_Belgorod_Oblast.svg",
			},
			{
				Code: "ru-bry",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2b/Flag_of_Bryansk_Oblast.svg",
			},
			{
				Code: "ru-bu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/68/Flag_of_Buryatia.svg",
			},
			{
				Code: "ru-ce",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/13/Flag_of_the_Chechen_Republic.svg",
			},
			{
				Code: "ru-che",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ce/Flag_of_Chelyabinsk_Oblast.svg",
			},
			{
				Code: "ru-chu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a0/Flag_of_Chukotka.svg",
			},
			{
				Code: "ru-cu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d7/Flag_of_Chuvashia.svg",
			},
			{
				Code: "ru-da",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c3/Flag_of_Dagestan.svg",
			},
			{
				Code: "ru-al",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/ff/Flag_of_Altai_Republic.svg",
			},
			{
				Code: "ru-in",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Flag_of_Ingushetia.svg",
			},
			{
				Code: "ru-irk",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/14/Flag_of_Irkutsk_Oblast.svg",
			},
			{
				Code: "ru-iva",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5b/Flag_of_Ivanovo_Oblast.svg",
			},
			{
				Code: "ru-kb",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d8/Flag_of_Kabardino-Balkaria.svg",
			},
			{
				Code: "ru-kc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/59/Flag_of_Karachay-Cherkessia.svg",
			},
			{
				Code: "ru-kda",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a7/Flag_of_Krasnodar_Krai.svg",
			},
			{
				Code: "ru-kem",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fd/Flag_of_Kemerovo_oblast.svg",
			},
			{
				Code: "ru-klu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ed/Flag_of_Kaluga_Oblast.svg",
			},
			{
				Code: "ru-kha",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_of_Khabarovsk_%28Khabarovsk_kray%29.png",
			},
			{
				Code: "ru-kr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/69/Flag_of_Karelia.svg",
			},
			{
				Code: "ru-kk",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ec/Flag_of_Khakassia.svg",
			},
			{
				Code: "ru-kl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9f/Flag_of_Kalmykia.svg",
			},
			{
				Code: "ru-khm",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/70/Flag_of_Yugra.svg",
			},
			{
				Code: "ru-kgd",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e4/Flag_of_Kaliningrad_Oblast.svg",
			},
			{
				Code: "ru-ko",
				Url:  "https://en.wikipedia.org/wiki/Komi_Republic#/media/File:Flag_of_Komi.svg",
			},
			{
				Code: "ru-kam",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a7/Flag_of_Kamchatka_Krai.svg",
			},
			{
				Code: "ru-krs",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/70/Flag_of_Kursk_Oblast.svg",
			},
			{
				Code: "ru-kos",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2c/Flag_of_Kostroma_Oblast.svg",
			},
			{
				Code: "ru-kgn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/92/Flag_of_Kurgan_Oblast.svg",
			},
			{
				Code: "ru-kir",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/41/Flag_of_Kirov_Oblast.svg",
			},
			{
				Code: "ru-kya",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2d/Flag_of_Krasnoyarsk_Krai.svg",
			},
			{
				Code: "ru-len",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/66/Flag_of_Leningrad_Oblast.svg",
			},
			{
				Code: "ru-lip",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/36/Flag_of_Lipetsk_Oblast.svg",
			},
			{
				Code: "ru-mow",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/03/Flag_of_Moscow%2C_Russia.svg",
			},
			{
				Code: "ru-me",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a7/Flag_of_Mari_El.svg",
			},
			{
				Code: "ru-mag",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/29/Flag_of_Magadan_Oblast.svg",
			},
			{
				Code: "ru-mur",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3b/Flag_of_Murmansk_Oblast.svg",
			},
			{
				Code: "ru-mo",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2f/Flag_of_Mordovia.svg",
			},
			{
				Code: "ru-mos",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/48/Flag_of_Moscow_oblast.svg",
			},
			{
				Code: "ru-ngr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/68/Flag_of_Novgorod_Oblast.svg",
			},
			{
				Code: "ru-nen",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/15/Flag_of_Nenets_Autonomous_District.svg",
			},
			{
				Code: "ru-se",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1c/Flag_of_North_Ossetia.svg",
			},
			{
				Code: "ru-nvs",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/39/Flag_of_Novosibirsk_oblast.svg",
			},
			{
				Code: "ru-niz",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/04/Flag_of_Nizhny_Novgorod_Region.svg",
			},
			{
				Code: "ru-ore",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d6/Flag_of_Orenburg_Oblast.svg",
			},
			{
				Code: "ru-orl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/05/Flag_of_Oryol_Oblast.svg",
			},
			{
				Code: "ru-oms",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/60/Flag_of_Omsk_Oblast.svg",
			},
			{
				Code: "ru-per",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b5/Flag_of_Perm_Krai.svg",
			},
			{
				Code: "ru-pri",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/38/Flag_of_Primorsky_Krai.svg",
			},
			{
				Code: "ru-psk",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ef/Flag_of_Pskov_Oblast.svg",
			},
			{
				Code: "ru-pnz",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/da/Flag_of_Penza_Oblast.svg",
			},
			{
				Code: "ru-ros",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d9/Flag_of_Rostov_Oblast.svg",
			},
			{
				Code: "ru-rya",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8c/Flag_of_Ryazan_Oblast.svg",
			},
			{
				Code: "ru-sam",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/13/Flag_of_Samara_Oblast.svg",
			},
			{
				Code: "ru-sa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/eb/Flag_of_Sakha.svg",
			},
			{
				Code: "ru-sak",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/57/Flag_of_Sakhalin_Oblast.svg",
			},
			{
				Code: "ru-smo",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bd/Flag_of_Smolensk_oblast.svg",
			},
			{
				Code: "ru-spe",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ca/Flag_of_Saint_Petersburg.svg",
			},
			{
				Code: "ru-sar",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f5/Flag_of_Saratov_Oblast.svg",
			},
			{
				Code: "ru-sta",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b4/Flag_of_Stavropol_Krai.svg",
			},
			{
				Code: "ru-sve",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ef/Flag_of_Sverdlovsk_Oblast.svg",
			},
			{
				Code: "ru-tam",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/39/Flag_of_Tambov_Oblast.svg",
			},
			{
				Code: "ru-tom",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/50/Flag_of_Tomsk_Oblast.svg",
			},
			{
				Code: "ru-tul",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/69/Flag_of_Tula_Oblast.svg",
			},
			{
				Code: "ru-ta",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/28/Flag_of_Tatarstan.svg",
			},
			{
				Code: "ru-ty",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/77/Flag_of_Tuva.svg",
			},
			{
				Code: "ru-tve",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/60/Flag_of_Tver_Oblast.svg",
			},
			{
				Code: "ru-tyu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4e/Flag_of_Tyumen_Oblast.svg",
			},
			{
				Code: "ru-ud",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c7/Flag_of_Udmurtia.svg",
			},
			{
				Code: "ru-uly",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d8/Flag_of_Ulyanovsk_Oblast.svg",
			},
			{
				Code: "ru-vgg",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/35/Flag_of_Volgograd_Oblast.svg",
			},
			{
				Code: "ru-vla",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/85/Flag_of_Vladimirskaya_Oblast.svg",
			},
			{
				Code: "ru-yan",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c7/Flag_of_Yamal-Nenets_Autonomous_District.svg",
			},
			{
				Code: "ru-vlg",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/af/Flag_of_Vologda_oblast.svg",
			},
			{
				Code: "ru-vor",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/64/Flag_of_Voronezh_Oblast.svg",
			},
			{
				Code: "ru-yar",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/ba/Flag_of_Yaroslavl_Oblast.svg",
			},
			{
				Code: "ru-yev",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0e/Flag_of_the_Jewish_Autonomous_Oblast.svg",
			},
			{
				Code: "ru-zab",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_of_Zabaykalsky_Krai.svg",
			},
		},
	},
	{
		Key:   "kr",
		Label: "🇰🇷 South Korea, Provinces",
		Entries: []FlagEntry{

			{
				Code: "kr-11",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/ca/Flag_of_Seoul.svg",
			},
			{
				Code: "kr-26",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b1/Flag_of_Busan.svg",
			},
			{
				Code: "kr-27",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Flag_of_Daegu.svg",
			},
			{
				Code: "kr-28",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0c/Flag_of_Incheon.svg",
			},
			{
				Code: "kr-29",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/59/Flag_of_Gwangju.svg",
			},
			{
				Code: "kr-30",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/53/Flag_of_Daejeon.svg",
			},
			{
				Code: "kr-31",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ef/Flag_of_Ulsan.svg",
			},
			{
				Code: "kr-41",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/80/Flag_of_Gyeonggi_Province_%282021%29.svg",
			},
			{
				Code: "kr-42",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fe/Flag_of_Gangwon_Province.svg",
			},
			{
				Code: "kr-43",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7d/Flag_of_North_Chungcheong_Province.svg",
			},
			{
				Code: "kr-44",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1c/Flag_of_South_Chungcheong_Province.svg",
			},
			{
				Code: "kr-45",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2a/Flag_of_North_Jeolla_Province.svg",
			},
			{
				Code: "kr-46",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/58/Flag_of_South_Jeolla_Province.svg",
			},
			{
				Code: "kr-47",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ef/Flag_of_North_Gyeongsang_Province.svg",
			},
			{
				Code: "kr-48",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ae/Flag_of_South_Gyeongsang_Province.svg",
			},
			{
				Code: "kr-49",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/35/Flag_of_Jeju.svg",
			},
			{
				Code: "kr-50",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/45/Flag_of_Sejong_City%2C_South_Korea.svg",
			},
		},
	},
	{
		Key:   "es",
		Label: "🇪🇸 Spain, Provinces",
		Entries: []FlagEntry{

			{
				Code: "es-a",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/75/Escut_de_la_Prov%C3%ADncia_d%27Alacant.svg",
			},
			{
				Code: "es-ab",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d6/Bandera_provincia_Albacete.svg",
			},
			{
				Code: "es-al",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/78/Flag_Almer%C3%ADa_Province.svg",
			},
			{
				Code: "es-av",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/54/Bandera_de_la_provincia_de_%C3%81vila.svg",
			},
			{
				Code: "es-b",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6e/Flag_of_Barcelona_province%28official%29.svg",
			},
			{
				Code: "es-ba",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8c/Provincia_de_Badajoz_-_Bandera.svg",
			},
			{
				Code: "es-bi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Bandera_de_Vizcaya.svg",
			},
			{
				Code: "es-bu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b5/Flag_Burgos_Province.svg",
			},
			{
				Code: "es-c",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ae/Coat_of_Arms_of_the_Province_of_Corunna.svg",
			},
			{
				Code: "es-ca",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d8/Flag_C%C3%A1diz_Province.svg",
			},
			{
				Code: "es-cc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/95/Flag_of_the_province_of_C%C3%A1ceres.svg",
			},
			{
				Code: "es-ce",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fd/Flag_Ceuta.svg",
			},
			{
				Code: "es-co",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/60/Provincia_de_C%C3%B3rdoba_-_Bandera.svg",
			},
			{
				Code: "es-cr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/37/Flag_Ciudad_Real_Province.svg",
			},
			{
				Code: "es-cs",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/79/Escut_de_la_Prov%C3%ADncia_de_Castell%C3%B3.svg",
			},
			{
				Code: "es-cu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2c/Flag_Cuenca_Province.svg",
			},
			{
				Code: "es-gc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/2d/Bandera_Provincial_de_Las_Palmas.svg",
			},
			{
				Code: "es-gi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/05/Flag_of_Girona_province_%28unofficial%29.svg",
			},
			{
				Code: "es-gr",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3a/Bandera_de_la_provincia_de_Granada_%28Espa%C3%B1a%29.svg",
			},
			{
				Code: "es-gu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8a/Flag_Guadalajara_Province.svg",
			},
			{
				Code: "es-h",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/24/Bandera_de_la_Provincia_De_Huelva.svg",
			},
			{
				Code: "es-hu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1f/Flag_of_Huesca_%28province%29.svg",
			},
			{
				Code: "es-j",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/87/Bandera_de_la_provincia_de_Ja%C3%A9n.svg",
			},
			{
				Code: "es-l",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/09/Bandera_de_la_provincia_de_L%C3%A9rida.svg",
			},
			{
				Code: "es-le",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8b/Bandera_de_Le%C3%B3n.svg",
			},
			{
				Code: "es-lo",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/db/Flag_of_La_Rioja_%28with_coat_of_arms%29.svg",
			},
			{
				Code: "es-lu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/fb/Flag_of_Lugo_province.svg",
			},
			{
				Code: "es-m",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9c/Flag_of_the_Community_of_Madrid.svg",
			},
			{
				Code: "es-ma",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/21/Bandera_Provincia_de_M%C3%A1laga.jpg",
			},
			{
				Code: "es-ml",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Melilla.svg",
			},
			{
				Code: "es-mu",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a5/Flag_of_the_Region_of_Murcia.svg",
			},
			{
				Code: "es-na",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/36/Bandera_de_Navarra.svg",
			},
			{
				Code: "es-o",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3e/Flag_of_Asturias.svg",
			},
			{
				Code: "es-or",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5e/Provincia_de_Ourense_-_Bandera.svg",
			},
			{
				Code: "es-p",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ea/Bandera_de_la_provincia_de_Palencia.svg",
			},
			{
				Code: "es-pm",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/7b/Flag_of_the_Balearic_Islands.svg",
			},
			{
				Code: "es-po",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a7/Flag_Pontevedra_Province.svg",
			},
			{
				Code: "es-s",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/0f/Flag_of_Cantabria_%28Official%29.svg",
			},
			{
				Code: "es-sa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Bandera_de_la_provincia_de_Salamanca.svg",
			},
			{
				Code: "es-se",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c4/Flag_of_Diputacion_de_Sevilla_Spain.svg",
			},
			{
				Code: "es-sg",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_Segovia_province.svg",
			},
			{
				Code: "es-so",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/79/Flag_Soria_province.svg",
			},
			{
				Code: "es-ss",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f5/Flag_of_Guip%C3%BAzcoa.svg",
			},
			{
				Code: "es-te",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Teruel_%28province%29.svg",
			},
			{
				Code: "es-tf",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/07/Bandera_Provincial_de_Santa_Cruz_de_Tenerife.svg",
			},
			{
				Code: "es-to",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bf/Bandera_de_la_provincia_de_Toledo.svg",
			},
			{
				Code: "es-t",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f6/Diputaci%C3%B3_copy.jpg",
			},
			{
				Code: "es-v",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1f/Escut_de_la_Prov%C3%ADncia_de_Val%C3%A8ncia.svg",
			},
			{
				Code: "es-va",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5c/Bandera_de_la_provincia_de_Valladolid.svg",
			},
			{
				Code: "es-vi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1f/Flag_of_%C3%81lava.svg",
			},
			{
				Code: "es-z",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e1/Flag_of_Zaragoza_province_%28with_coat_of_arms%29.svg",
			},
			{
				Code: "es-za",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5a/Coat_of_Arms_of_Zamora_Province.svg",
			},
		},
	},
	{
		Key:   "ua",
		Label: "🇺🇦 Ukraine, Oblasts",
		Entries: []FlagEntry{

			{
				Code: "ua-43",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/aa/Flag_of_Crimea.svg",
			},
			{
				Code: "ua-71",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/85/Flag_of_Cherkasy_Oblast.svg",
			},
			{
				Code: "ua-74",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Chernihiv_Oblast.svg",
			},
			{
				Code: "ua-77",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a6/Flag_of_Chernivtsi_Oblast.svg",
			},
			{
				Code: "ua-12",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/cf/Flag_of_Dnipropetrovsk_Oblast.svg",
			},
			{
				Code: "ua-14",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5d/Flag_of_Donetsk_Oblast.svg",
			},
			{
				Code: "ua-26",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/65/Flag_of_Ivano-Frankivsk_Oblast2.svg",
			},
			{
				Code: "ua-63",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/41/Flag_of_Kharkiv_Oblast.svg",
			},
			{
				Code: "ua-65",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c3/Flag_of_Kherson_Oblast.svg",
			},
			{
				Code: "ua-68",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/37/Flag_of_Khmelnytskyi_Oblast.svg",
			},
			{
				Code: "ua-32",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b2/Flag_of_Kyiv_Oblast.svg",
			},
			{
				Code: "ua-35",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/09/Flag_of_Kirovohrad_Oblast.svg",
			},
			{
				Code: "ua-09",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/7/77/Flag_of_Luhansk_Oblast.svg",
			},
			{
				Code: "ua-46",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/82/Flag_of_Lviv_Oblast.svg",
			},
			{
				Code: "ua-48",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e1/Flag_of_Mykolaiv_Oblast.svg",
			},
			{
				Code: "ua-51",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/ba/Flag_of_Odesa_Oblast.svg",
			},
			{
				Code: "ua-53",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d3/Flag_of_Poltava_Oblast.svg",
			},
			{
				Code: "ua-56",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3c/Flag_of_Rivne_Oblast.svg",
			},
			{
				Code: "ua-59",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/21/Flag_of_Sumy_Oblast.svg",
			},
			{
				Code: "ua-61",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/dc/Flag_of_Ternopil_Oblast.svg",
			},
			{
				Code: "ua-05",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/03/Flag_of_Vinnytsia_Oblast.svg",
			},
			{
				Code: "ua-07",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/00/Flag_of_Volyn_Oblast.svg",
			},
			{
				Code: "ua-21",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/3c/Flag_of_Transcarpathian_Oblast.svg",
			},
			{
				Code: "ua-23",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b8/Flag_of_Zaporizhia_Oblast.svg",
			},
			{
				Code: "ua-18",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/69/Flag_of_Zhytomyr_Oblast.svg",
			},
		},
	},
	{
		Key:   "us",
		Label: "🇺🇸 US, States",
		Entries: []FlagEntry{
			{
				Code: "us-al",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5c/Flag_of_Alabama.svg",
			},
			{
				Code: "us-ak",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e6/Flag_of_Alaska.svg",
			},
			{
				Code: "us-az",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9d/Flag_of_Arizona.svg",
			},
			{
				Code: "us-ar",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9d/Flag_of_Arkansas.svg",
			},
			{
				Code: "us-ca",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/01/Flag_of_California.svg",
			},
			{
				Code: "us-co",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/21/Flag_of_Colorado_designed_by_Andrew_Carlisle_Carson.svg",
			},
			{
				Code: "us-ct",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/96/Flag_of_Connecticut.svg",
			},
			{
				Code: "us-de",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c6/Flag_of_Delaware.svg",
			},
			{
				Code: "us-fl",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Florida.svg",
			},
			{
				Code: "us-ga",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/54/Flag_of_Georgia_(U.S._state).svg",
			},
			{
				Code: "us-hi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ef/Flag_of_Hawaii.svg",
			},
			{
				Code: "us-id",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a4/Flag_of_Idaho.svg",
			},
			{
				Code: "us-il",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/0/01/Flag_of_Illinois.svg",
			},
			{
				Code: "us-in",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/ac/Flag_of_Indiana.svg",
			},
			{
				Code: "us-ia",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/aa/Flag_of_Iowa.svg",
			},
			{
				Code: "us-ks",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/da/Flag_of_Kansas.svg",
			},
			{
				Code: "us-ky",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/8/8d/Flag_of_Kentucky.svg",
			},
			{
				Code: "us-la",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/e0/Flag_of_Louisiana.svg",
			},
			{
				Code: "us-me",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/3/35/Flag_of_Maine.svg",
			},
			{
				Code: "us-md",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/a/a0/Flag_of_Maryland.svg",
			},
			{
				Code: "us-ma",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f2/Flag_of_Massachusetts.svg",
			},
			{
				Code: "us-mi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b5/Flag_of_Michigan.svg",
			},
			{
				Code: "us-mn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b9/Flag_of_Minnesota.svg",
			},
			{
				Code: "us-ms",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/42/Flag_of_Mississippi.svg",
			},
			{
				Code: "us-mo",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/5a/Flag_of_Missouri.svg",
			},
			{
				Code: "us-mt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/cb/Flag_of_Montana.svg",
			},
			{
				Code: "us-ne",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4d/Flag_of_Nebraska.svg",
			},
			{
				Code: "us-nv",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f1/Flag_of_Nevada.svg",
			},
			{
				Code: "us-nh",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/28/Flag_of_New_Hampshire.svg",
			},
			{
				Code: "us-nj",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/92/Flag_of_New_Jersey.svg",
			},
			{
				Code: "us-nm",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/c/c3/Flag_of_New_Mexico.svg",
			},
			{
				Code: "us-ny",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1a/Flag_of_New_York.svg",
			},
			{
				Code: "us-nc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bb/Flag_of_North_Carolina.svg",
			},
			{
				Code: "us-nd",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/e/ee/Flag_of_North_Dakota.svg",
			},
			{
				Code: "us-oh",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/4c/Flag_of_Ohio.svg",
			},
			{
				Code: "us-ok",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/6e/Flag_of_Oklahoma.svg",
			},
			{
				Code: "us-or",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/b9/Flag_of_Oregon.svg",
			},
			{
				Code: "us-pa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Pennsylvania.svg",
			},
			{
				Code: "us-ri",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f3/Flag_of_Rhode_Island.svg",
			},
			{
				Code: "us-sc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/6/69/Flag_of_South_Carolina.svg",
			},
			{
				Code: "us-sd",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/1/1a/Flag_of_South_Dakota.svg",
			},
			{
				Code: "us-tn",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/9/9e/Flag_of_Tennessee.svg",
			},
			{
				Code: "us-tx",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f7/Flag_of_Texas.svg",
			},
			{
				Code: "us-ut",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/f/f6/Flag_of_Utah.svg",
			},
			{
				Code: "us-vt",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/49/Flag_of_Vermont.svg",
			},
			{
				Code: "us-va",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/4/47/Flag_of_Virginia.svg",
			},
			{
				Code: "us-wa",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/5/54/Flag_of_Washington.svg",
			},
			{
				Code: "us-dc",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/d/d4/Flag_of_the_District_of_Columbia.svg",
			},
			{
				Code: "us-wv",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/22/Flag_of_West_Virginia.svg",
			},
			{
				Code: "us-wi",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/2/22/Flag_of_Wisconsin.svg",
			},
			{
				Code: "us-wy",
				Url:  "https://upload.wikimedia.org/wikipedia/commons/b/bc/Flag_of_Wyoming.svg",
			},
		},
	},
}
