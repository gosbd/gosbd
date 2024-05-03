package lang_test

import (
	"testing"

	"github.com/gosbd/gosbd"
)

func Test_English(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1) Simple period to end sentence",
			args: args{
				text: "Hello World. My name is Jonas.",
			},
			want: []string{"Hello World.", "My name is Jonas."},
		},
		{
			name: "2) Question mark to end sentence",
			args: args{
				text: "What is your name? My name is Jonas.",
			},
			want: []string{"What is your name?", "My name is Jonas."},
		},
		{
			name: "3) Exclamation point to end sentence",
			args: args{
				text: "There it is! I found it.",
			},
			want: []string{"There it is!", "I found it."},
		},
		{
			name: "4) One letter upper case abbreviations",
			args: args{
				text: "My name is Jonas E. Smith.",
			},
			want: []string{"My name is Jonas E. Smith."},
		},
		{
			name: "5) One letter lower case abbreviations",
			args: args{
				text: "Please turn to p. 55.",
			},
			want: []string{"Please turn to p. 55."},
		},
		{
			name: "6) Two letter lower case abbreviations in the middle of a sentence",
			args: args{
				text: "Were Jane and co. at the party?",
			},
			want: []string{"Were Jane and co. at the party?"},
		},
		{
			name: "7) Two letter upper case abbreviations in the middle of a sentence",
			args: args{
				text: "They closed the deal with Pitt, Briggs & Co. at noon.",
			},
			want: []string{"They closed the deal with Pitt, Briggs & Co. at noon."},
		},
		{
			name: "8) Two letter lower case abbreviations at the end of a sentence",
			args: args{
				text: "Let's ask Jane and co. They should know.",
			},
			want: []string{"Let's ask Jane and co.", "They should know."},
		},
		{
			name: "9) Two letter upper case abbreviations at the end of a sentence",
			args: args{
				text: "They closed the deal with Pitt, Briggs & Co. It closed yesterday.",
			},
			want: []string{"They closed the deal with Pitt, Briggs & Co.", "It closed yesterday."},
		},
		{
			name: "10) Two letter (prepositive) abbreviations",
			args: args{
				text: "I can see Mt. Fuji from here.",
			},
			want: []string{"I can see Mt. Fuji from here."},
		},
		{
			name: "11) Two letter (prepositive & postpositive) abbreviations",
			args: args{
				text: "St. Michael's Church is on 5th st. near the light.",
			},
			want: []string{"St. Michael's Church is on 5th st. near the light."},
		},
		{
			name: "12) Possesive two letter abbreviations",
			args: args{
				text: "That is JFK Jr.'s book.",
			},
			want: []string{"That is JFK Jr.'s book."},
		},
		{
			name: "13) Multi-period abbreviations in the middle of a sentence",
			args: args{
				text: "I visited the U.S.A. last year.",
			},
			want: []string{"I visited the U.S.A. last year."},
		},
		{
			name: "14) Multi-period abbreviations at the end of a sentence",
			args: args{
				text: "I live in the E.U. How about you?",
			},
			want: []string{"I live in the E.U.", "How about you?"},
		},
		{
			name: "15) U.S. as sentence boundary",
			args: args{
				text: "I live in the U.S. How about you?",
			},
			want: []string{"I live in the U.S.", "How about you?"},
		},
		{
			name: "16) U.S. as non sentence boundary with next word capitalized",
			args: args{
				text: "I work for the U.S. Government in Virginia.",
			},
			want: []string{"I work for the U.S. Government in Virginia."},
		},
		{
			name: "17) U.S. as non sentence boundary",
			args: args{
				text: "I have lived in the U.S. for 20 years.",
			},
			want: []string{"I have lived in the U.S. for 20 years."},
		},
		/*
			{
				// Most difficult sentence to crack
				name: "18) A.M. / P.M. as non sentence boundary and sentence boundary",
				args: args{
					text: "At 5 a.m. Mr. Smith went to the bank. He left the bank at 6 P.M. Mr. Smith then went to the store.",
				},
				want: []string{"At 5 a.m. Mr. Smith went to the bank.", "He left the bank at 6 P.M.", "Mr. Smith then went to the store."},
			},
		*/
		{
			name: "19) Number as non sentence boundary",
			args: args{
				text: "She has $100.00 in her bag.",
			},
			want: []string{"She has $100.00 in her bag."},
		},
		{
			name: "20) Number as sentence boundary",
			args: args{
				text: "She has $100.00. It is in her bag.",
			},
			want: []string{"She has $100.00.", "It is in her bag."},
		},
		{
			name: "21) Parenthetical inside sentence",
			args: args{
				text: "He teaches science (He previously worked for 5 years as an engineer.) at the local University.",
			},
			want: []string{"He teaches science (He previously worked for 5 years as an engineer.) at the local University."},
		},
		{
			name: "22) Email addresses",
			args: args{
				text: "Her email is Jane.Doe@example.com. I sent her an email.",
			},
			want: []string{
				"Her email is Jane.Doe@example.com.", "I sent her an email.",
			},
		},
		{
			name: "23) Web addresses",
			args: args{
				text: "The site is: https://www.example.50.com/new-site/awesome_content.html. Please check it out.",
			},
			want: []string{
				"The site is: https://www.example.50.com/new-site/awesome_content.html.", "Please check it out.",
			},
		},
		{
			name: "24) Single quotations inside sentence",
			args: args{
				text: "She turned to him, 'This is great.' she said.",
			},
			want: []string{
				"She turned to him, 'This is great.' she said.",
			},
		},
		{
			name: "25) Double quotations inside sentence",
			args: args{
				text: `She turned to him, "This is great." she said.`,
			},
			want: []string{
				`She turned to him, "This is great." she said.`,
			},
		},
		{
			name: "26) Double quotations at the end of a sentence",
			args: args{
				text: `She turned to him, "This is great." She held the book out to show him.`,
			},
			want: []string{
				`She turned to him, "This is great."`, "She held the book out to show him.",
			},
		},
		{
			name: "27) Double punctuation (exclamation point)",
			args: args{
				text: "Hello!! Long time no see.",
			},
			want: []string{
				"Hello!!", "Long time no see.",
			},
		},
		{
			name: "28) Double punctuation (question mark)",
			args: args{
				text: "Hello?? Who is there?",
			},
			want: []string{
				"Hello??", "Who is there?",
			},
		},
		{
			name: "29) Double punctuation (exclamation point / question mark)",
			args: args{
				text: "Hello!? Is that you?",
			},
			want: []string{
				"Hello!?", "Is that you?",
			},
		},
		{
			name: "30) Double punctuation (question mark / exclamation point)",
			args: args{
				text: "Hello?! Is that you?",
			},
			want: []string{
				"Hello?!", "Is that you?",
			},
		},
		{
			name: "31) List (period followed by parens and no period to end item)",
			args: args{
				text: "1.) The first item 2.) The second item",
			},
			want: []string{
				"1.) The first item", "2.) The second item",
			},
		},
		{
			name: "32) List (period followed by parens and period to end item)",
			args: args{
				text: "1.) The first item. 2.) The second item.",
			},
			want: []string{
				"1.) The first item.", "2.) The second item.",
			},
		},
		{
			name: "33) List (parens and no period to end item)",
			args: args{
				text: "1) The first item 2) The second item",
			},
			want: []string{
				"1) The first item", "2) The second item",
			},
		},
		{
			name: "34) List (parens and period to end item)",
			args: args{
				text: "1) The first item. 2) The second item.",
			},
			want: []string{
				"1) The first item.", "2) The second item.",
			},
		},
		{
			name: "35) List (period to mark list and no period to end item)",
			args: args{
				text: "1. The first item 2. The second item",
			},
			want: []string{
				"1. The first item", "2. The second item",
			},
		},
		{
			name: "36) List (period to mark list and period to end item)",
			args: args{
				text: "1. The first item. 2. The second item.",
			},
			want: []string{
				"1. The first item.", "2. The second item.",
			},
		},
		{
			name: "37) List with bullet",
			args: args{
				text: "• 9. The first item • 10. The second item",
			},
			want: []string{
				"• 9. The first item", "• 10. The second item",
			},
		},
		{
			name: "38) List with hypthen",
			args: args{
				text: "⁃9. The first item ⁃10. The second item",
			},
			want: []string{
				"⁃9. The first item", "⁃10. The second item",
			},
		},
		{
			name: "39) Alphabetical list",
			args: args{
				text: "a. The first item b. The second item c. The third list item",
			},
			want: []string{
				"a. The first item", "b. The second item", "c. The third list item",
			},
		},
		{
			name: "40) Geo Coordinates",
			args: args{
				text: "You can find it at N°. 1026.253.553. That is where the treasure is.",
			},
			want: []string{
				"You can find it at N°. 1026.253.553.", "That is where the treasure is.",
			},
		},
		{
			name: "41) Named entities with an exclamation point",
			args: args{
				text: "She works at Yahoo! in the accounting department.",
			},
			want: []string{
				"She works at Yahoo! in the accounting department.",
			},
		},
		{
			name: "42) I as a sentence boundary and I as an abbreviation",
			args: args{
				text: "We make a good team, you and I. Did you see Albert I. Jones yesterday?",
			},
			want: []string{
				"We make a good team, you and I.", "Did you see Albert I. Jones yesterday?",
			},
		},
		{
			name: "43) Ellipsis at end of quotation",
			args: args{
				text: `Thoreau argues that by simplifying one’s life, “the laws of the universe will appear less complex. . . .”`,
			},
			want: []string{
				`Thoreau argues that by simplifying one’s life, “the laws of the universe will appear less complex. . . .”`,
			},
		},
		{
			name: "44) Ellipsis with square brackets",
			args: args{
				text: `"Bohr [...] used the analogy of parallel stairways [...]" (Smith 55).`,
			},
			want: []string{
				`"Bohr [...] used the analogy of parallel stairways [...]" (Smith 55).`,
			},
		},
		{
			name: "45) Ellipsis as sentence boundary (standard ellipsis rules)",
			args: args{
				text: "If words are left off at the end of a sentence, and that is all that is omitted, indicate the omission with ellipsis marks (preceded and followed by a space) and then indicate the end of the sentence with a period . . . . Next sentence.",
			},
			want: []string{
				"If words are left off at the end of a sentence, and that is all that is omitted, indicate the omission with ellipsis marks (preceded and followed by a space) and then indicate the end of the sentence with a period . . . .",
				"Next sentence.",
			},
		},
		{
			name: "46) Ellipsis as sentence boundary (standard ellipsis rules)",
			args: args{
				text: "I never meant that.... She left the store.",
			},
			want: []string{
				"I never meant that....", "She left the store.",
			},
		},
		{
			name: "47) Ellipsis as non sentence boundary",
			args: args{
				text: "I wasn’t really ... well, what I mean...see . . . what I'm saying, the thing is . . . I didn’t mean it.",
			},
			want: []string{
				"I wasn’t really ... well, what I mean...see . . . what I'm saying, the thing is . . . I didn’t mean it.",
			},
		},
		{
			name: "48) 4-dot ellipsis",
			args: args{
				text: "One further habit which was somewhat weakened . . . was that of combining words into self-interpreting compounds. . . . The practice was not abandoned. . . .",
			},
			want: []string{
				"One further habit which was somewhat weakened . . . was that of combining words into self-interpreting compounds.",
				". . . The practice was not abandoned. . . .",
			},
		},
		{
			name: "Bugfix #12",
			args: args{
				text: `Candidates tied to Tehreek-e-Insaf (PTI), the party of Imran Khan, won the most seats in Pakistan’s general election, despite a de facto ban on their campaign. Mr Khan is in prison on multiple charges, which he says are politically motivated. The Pakistan Muslim League-Nawaz (PML-N), which was widely expected to win, came second. PML-N is the party of Nawaz Sharif, Mr Khan’s arch-rival. It will form a coalition government with the Pakistan Peoples Party, which came third. Mr Khan’s supporters said the election had been rigged, which the PML-N denied. The head of the army claimed the poll had been “free and unhindered”.`,
			},
			want: []string{
				"Candidates tied to Tehreek-e-Insaf (PTI), the party of Imran Khan, won the most seats in Pakistan’s general election, despite a de facto ban on their campaign.",
				"Mr Khan is in prison on multiple charges, which he says are politically motivated.",
				"The Pakistan Muslim League-Nawaz (PML-N), which was widely expected to win, came second.",
				"PML-N is the party of Nawaz Sharif, Mr Khan’s arch-rival.",
				"It will form a coalition government with the Pakistan Peoples Party, which came third.",
				"Mr Khan’s supporters said the election had been rigged, which the PML-N denied.",
				"The head of the army claimed the poll had been “free and unhindered”.",
			},
		},
		{
			name: "Regression #12",
			args: args{
				text: "Mix it, put it in the oven, and -- voila! -- you have cake. Some can be -- if I may say so? -- a bit questionable.",
			},
			want: []string{
				"Mix it, put it in the oven, and -- voila! -- you have cake.",
				"Some can be -- if I may say so? -- a bit questionable.",
			},
		},
		{
			name: "Issue #14",
			args: args{
				text: `The Academy Award for Best Production Design recognizes achievement for art direction in film. The category's original name was Best Art Direction, but was changed to its current name in 2012 for the 85th Academy Awards.[1] This change resulted from the Art Directors' bggranch of the Academy of Motion Picture Arts and Sciences (AMPAS) being renamed the Designers' branch. Since 1947, the award is shared with the set decorators. It is awarded to the best interior design in a film.[2] The films below are listed with their production year (for example, the 2000 Academy Award for Best Art Direction is given to a film from 1999). In the lists below, the winner of the award for each year is shown first, followed by the other nominees in alphabetical order.`,
			},
			want: []string{
				"The Academy Award for Best Production Design recognizes achievement for art direction in film.",
				"The category's original name was Best Art Direction, but was changed to its current name in 2012 for the 85th Academy Awards.[1]",
				"This change resulted from the Art Directors' bggranch of the Academy of Motion Picture Arts and Sciences (AMPAS) being renamed the Designers' branch.",
				"Since 1947, the award is shared with the set decorators.",
				"It is awarded to the best interior design in a film.[2]",
				"The films below are listed with their production year (for example, the 2000 Academy Award for Best Art Direction is given to a film from 1999).",
				"In the lists below, the winner of the award for each year is shown first, followed by the other nominees in alphabetical order.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sg := gosbd.NewSegmenter("en")
			got := sg.Segment(tt.args.text)
			for i, v := range got {
				if i >= len(tt.want) {
					break
				}
				if v != tt.want[i] {
					t.Errorf("Segmenter.Segment() = %#v, want %#v", v, tt.want[i])
				}
			}
			if len(got) != len(tt.want) {
				t.Errorf("Segmenter.Segment() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
