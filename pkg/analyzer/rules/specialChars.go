package rules

import (
	"go/token"
	"unicode"
	"unicode/utf8"

	"github.com/m4kson/loglint/pkg/analyzer/detector"
	"golang.org/x/tools/go/analysis"
)

type NoSpecialCharsRule struct{}

func (r *NoSpecialCharsRule) Name() string { return "no-special-chars" }

var emojiRangeTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x2194, 0x2199, 1}, // arrows
		{0x21A9, 0x21AA, 1}, // arrows with hook
		{0x231A, 0x231B, 1}, // watch, hourglass
		{0x2328, 0x2328, 1}, // keyboard
		{0x23CF, 0x23CF, 1}, // eject
		{0x23E9, 0x23F3, 1}, // various clocks / media controls
		{0x23F8, 0x23FA, 1}, // pause / stop / record
		{0x24C2, 0x24C2, 1}, // circled M
		{0x25AA, 0x25AB, 1}, // small squares
		{0x25B6, 0x25B6, 1}, // play button
		{0x25C0, 0x25C0, 1}, // reverse button
		{0x25FB, 0x25FE, 1}, // medium squares
		{0x2600, 0x2604, 1}, // sun, cloud, umbrella…
		{0x260E, 0x260E, 1}, // telephone
		{0x2611, 0x2611, 1}, // ballot box
		{0x2614, 0x2615, 1}, // umbrella with rain, hot beverage
		{0x2618, 0x2618, 1}, // shamrock
		{0x261D, 0x261D, 1}, // index finger
		{0x2620, 0x2620, 1}, // skull and crossbones
		{0x2622, 0x2623, 1}, // radioactive, biohazard
		{0x2626, 0x2626, 1}, // orthodox cross
		{0x262A, 0x262A, 1}, // star and crescent
		{0x262E, 0x262F, 1}, // peace / yin yang
		{0x2638, 0x263A, 1}, // wheel of dharma, smiley
		{0x2640, 0x2640, 1}, // female sign
		{0x2642, 0x2642, 1}, // male sign
		{0x2648, 0x2653, 1}, // zodiac signs
		{0x265F, 0x2660, 1}, // chess pawn, spade
		{0x2663, 0x2663, 1}, // club
		{0x2665, 0x2666, 1}, // heart, diamond
		{0x2668, 0x2668, 1}, // hot springs
		{0x267B, 0x267B, 1}, // recycling
		{0x267E, 0x267F, 1}, // infinity, wheelchair
		{0x2692, 0x2697, 1}, // tools
		{0x2699, 0x2699, 1}, // gear
		{0x269B, 0x269C, 1}, // atom, fleur-de-lis
		{0x26A0, 0x26A1, 1}, // warning, lightning
		{0x26AA, 0x26AB, 1}, // circles
		{0x26B0, 0x26B1, 1}, // coffin, urn
		{0x26BD, 0x26BE, 1}, // soccer, baseball
		{0x26C4, 0x26C5, 1}, // snowman, sun behind cloud
		{0x26CE, 0x26CF, 1}, // ophiuchus, pick
		{0x26D1, 0x26D1, 1}, // helmet
		{0x26D3, 0x26D4, 1}, // chains, no entry
		{0x26E9, 0x26EA, 1}, // shinto shrine, church
		{0x26F0, 0x26F5, 1}, // mountain, sailboat
		{0x26F7, 0x26FA, 1}, // skier, tent
		{0x26FD, 0x26FD, 1}, // fuel pump
		{0x2702, 0x2702, 1}, // scissors
		{0x2705, 0x2705, 1}, // check mark
		{0x2708, 0x270D, 1}, // airplane … writing hand
		{0x270F, 0x270F, 1}, // pencil
		{0x2712, 0x2712, 1}, // black nib
		{0x2714, 0x2714, 1}, // heavy check
		{0x2716, 0x2716, 1}, // heavy multiplication x
		{0x271D, 0x271D, 1}, // latin cross
		{0x2721, 0x2721, 1}, // star of david
		{0x2728, 0x2728, 1}, // sparkles
		{0x2733, 0x2734, 1}, // eight-pointed stars
		{0x2744, 0x2744, 1}, // snowflake
		{0x2747, 0x2747, 1}, // sparkle
		{0x274C, 0x274C, 1}, // cross mark
		{0x274E, 0x274E, 1}, // cross mark button
		{0x2753, 0x2755, 1}, // question marks
		{0x2757, 0x2757, 1}, // exclamation
		{0x2763, 0x2764, 1}, // hearts
		{0x2795, 0x2797, 1}, // plus, minus, division
		{0x27A1, 0x27A1, 1}, // right arrow
		{0x27B0, 0x27B0, 1}, // curly loop
		{0x27BF, 0x27BF, 1}, // double curly loop
		{0x2934, 0x2935, 1}, // arrows
		{0x2B05, 0x2B07, 1}, // arrows
		{0x2B1B, 0x2B1C, 1}, // squares
		{0x2B50, 0x2B50, 1}, // star
		{0x2B55, 0x2B55, 1}, // circle
		{0x3030, 0x3030, 1}, // wavy dash
		{0x303D, 0x303D, 1}, // part alternation mark
		{0x3297, 0x3297, 1}, // circled ideograph congratulation
		{0x3299, 0x3299, 1}, // circled ideograph secret
	},
	R32: []unicode.Range32{
		{0x1F004, 0x1F004, 1}, // mahjong red dragon
		{0x1F0CF, 0x1F0CF, 1}, // joker
		{0x1F170, 0x1F171, 1}, // blood type
		{0x1F17E, 0x1F17F, 1}, // parking
		{0x1F18E, 0x1F18E, 1}, // AB button
		{0x1F191, 0x1F19A, 1}, // squared CL … squared VS
		{0x1F1E0, 0x1F1FF, 1}, // regional indicator symbols (flags)
		{0x1F201, 0x1F202, 1}, // squared CJK
		{0x1F21A, 0x1F21A, 1}, // squared CJK free
		{0x1F22F, 0x1F22F, 1}, // squared CJK reserved
		{0x1F232, 0x1F23A, 1}, // squared CJK ideographs
		{0x1F250, 0x1F251, 1}, // circled ideographs
		{0x1F300, 0x1F321, 1}, // misc symbols and pictographs
		{0x1F324, 0x1F393, 1}, // misc symbols and pictographs cont.
		{0x1F396, 0x1F397, 1}, // military medal, reminder ribbon
		{0x1F399, 0x1F39B, 1}, // studio microphone…
		{0x1F39E, 0x1F3F0, 1}, // film frames … european castle
		{0x1F3F3, 0x1F3F5, 1}, // flags
		{0x1F3F7, 0x1F4FD, 1}, // label … film projector
		{0x1F4FF, 0x1F53D, 1}, // prayer beads … down button
		{0x1F549, 0x1F54E, 1}, // om … menorah
		{0x1F550, 0x1F567, 1}, // clocks
		{0x1F56F, 0x1F570, 1}, // candle, mantelpiece clock
		{0x1F573, 0x1F57A, 1}, // hole … man dancing
		{0x1F587, 0x1F587, 1}, // linked paperclips
		{0x1F58A, 0x1F58D, 1}, // pens
		{0x1F590, 0x1F590, 1}, // hand with fingers splayed
		{0x1F595, 0x1F596, 1}, // middle finger, vulcan salute
		{0x1F5A4, 0x1F5A5, 1}, // black heart, desktop computer
		{0x1F5A8, 0x1F5A8, 1}, // printer
		{0x1F5B1, 0x1F5B2, 1}, // mouse buttons
		{0x1F5BC, 0x1F5BC, 1}, // frame with picture
		{0x1F5C2, 0x1F5C4, 1}, // card index dividers…
		{0x1F5D1, 0x1F5D3, 1}, // wastebasket…
		{0x1F5DC, 0x1F5DE, 1}, // compression…
		{0x1F5E1, 0x1F5E1, 1}, // dagger
		{0x1F5E3, 0x1F5E3, 1}, // speaking head
		{0x1F5E8, 0x1F5E8, 1}, // left speech bubble
		{0x1F5EF, 0x1F5EF, 1}, // right anger bubble
		{0x1F5F3, 0x1F5F3, 1}, // ballot box with ballot
		{0x1F5FA, 0x1F64F, 1}, // world map … folded hands
		{0x1F680, 0x1F6C5, 1}, // rocket … left luggage
		{0x1F6CB, 0x1F6D2, 1}, // couch … shopping cart
		{0x1F6D5, 0x1F6D7, 1}, // hindu temple…
		{0x1F6E0, 0x1F6E5, 1}, // hammer and wrench…
		{0x1F6E9, 0x1F6E9, 1}, // small airplane
		{0x1F6EB, 0x1F6EC, 1}, // airplane departure/arrival
		{0x1F6F0, 0x1F6F0, 1}, // satellite
		{0x1F6F3, 0x1F6FC, 1}, // passenger ship…
		{0x1F7E0, 0x1F7EB, 1}, // coloured circles and squares
		{0x1F90C, 0x1F93A, 1}, // pinched fingers…
		{0x1F93C, 0x1F945, 1}, // people wrestling…
		{0x1F947, 0x1F9FF, 1}, // 1st place medal…
		{0x1FA00, 0x1FA53, 1}, // chess pieces…
		{0x1FA60, 0x1FA6D, 1}, // games
		{0x1FA70, 0x1FA74, 1}, // ballet shoes…
		{0x1FA78, 0x1FA7A, 1}, // drop of blood…
		{0x1FA80, 0x1FA86, 1}, // yo-yo…
		{0x1FA90, 0x1FAA8, 1}, // ringed planet…
		{0x1FAB0, 0x1FAB6, 1}, // fly…
		{0x1FAC0, 0x1FAC2, 1}, // anatomical heart…
		{0x1FAD0, 0x1FAD6, 1}, // blueberries…
	},
}

func (r *NoSpecialCharsRule) Check(pass *analysis.Pass, call detector.Call) {
	msg := call.MsgValue
	if msg == "" {
		return
	}

	byteOffset := 0

	for _, ch := range msg {
		if ch == utf8.RuneError {
			break
		}

		reason := classifyBadRune(ch)
		if reason != "" {
			offendingPos := call.MsgPos + 1 + token.Pos(byteOffset)
			pass.Reportf(
				offendingPos,
				"log message must not contain %s, found %q",
				reason,
				string(ch),
			)
			return
		}

		size := utf8.RuneLen(ch)
		if size < 0 {
			size = 1
		}
		byteOffset += size
	}
}

func classifyBadRune(ch rune) string {
	if ch >= 0x20 && ch <= 0x7E {
		if ch != ' ' && ch != '%' && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return "special characters"
		}
		return ""
	}

	if unicode.Is(emojiRangeTable, ch) {
		return "emoji"
	}

	return ""
}
