package gox12

import (
	//"fmt"
	"strings"
	"testing"
)

func TestSegmentParseFormatIdentity(t *testing.T) {
	paths := [...]string{
		"/2000A/2000B/2300/2400/SV2",
		"/2000A/2000B/2300/2400/SV201",
		"/2000A/2000B/2300/2400/SV2[421]01",
	}
	for _, p := range paths {
		//fmt.Println(p)
		fpath, err := NewX12Path(p)
		if err != nil {
			t.Errorf("Didn't get a value for [%s]", p)
		}
		//fmt.Println(fpath)
		actual := fpath.String()
		if actual != p {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", p, actual)
		}
	}
}

func TestRefDes(t *testing.T) {
	var tests = []struct {
		spath     string
		seg_id    string
		qual      string
		eleidx    int
		subeleidx int
	}{
		{"TST", "TST", "", 0, 0},
		{"TST02", "TST", "", 2, 0},
		{"TST03-2", "TST", "", 3, 2},
		{"TST[AA]02", "TST", "AA", 2, 0},
		{"TST[1B5]03-1", "TST", "1B5", 3, 1},
		{"03", "", "", 3, 0},
		{"03-2", "", "", 3, 2},
		{"N102", "N1", "", 2, 0},
		{"N102-5", "N1", "", 2, 5},
		{"N1[AZR]02", "N1", "AZR", 2, 0},
		{"N1[372]02-5", "N1", "372", 2, 5},
	}
	for _, tt := range tests {
		actual, err := NewX12Path(tt.spath)
		if err != nil {
			t.Errorf("Didn't get a value for [%s]", actual)
		}
		if actual.SegmentId != tt.seg_id {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.seg_id, actual.SegmentId)
		}
		if actual.IdValue != tt.qual {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.qual, actual.IdValue)
		}
		if actual.ElementIdx != tt.eleidx {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.eleidx, actual.ElementIdx)
		}
		if actual.SubelementIdx != tt.subeleidx {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.subeleidx, actual.SubelementIdx)
		}
		if len(actual.Loops) != 0 {
			t.Errorf("Loops is not empty")
		}
		path := actual.String()
		if path != tt.spath {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.spath, path)
		}
	}
}

func TestRelativePath(t *testing.T) {
	var tests = []struct {
		spath     string
		seg_id    string
		qual      string
		eleidx    int
		subeleidx int
		loops     []string
	}{
		{"AAA/TST", "TST", "", 0, 0, []string{"AAA"}},
		{"B1000/TST02", "TST", "", 2, 0, []string{"B1000"}},
		{"1000B/TST03-2", "TST", "", 3, 2, []string{"1000B"}},
		{"1000A/1000B/TST[AA]02", "TST", "AA", 2, 0, []string{"1000A", "1000B"}},
		{"AA/BB/CC/TST[1B5]03-1", "TST", "1B5", 3, 1, []string{"AA", "BB", "CC"}},
		{"DDD/E1000/N102", "N1", "", 2, 0, []string{"DDD", "E1000"}},
		{"E1000/D322/N102-5", "N1", "", 2, 5, []string{"E1000", "D322"}},
		{"BB/CC/N1[AZR]02", "N1", "AZR", 2, 0, []string{"BB", "CC"}},
		{"BB/CC/N1[372]02-5", "N1", "372", 2, 5, []string{"BB", "CC"}},
	}
	for _, tt := range tests {
		actual, err := NewX12Path(tt.spath)
		if err != nil {
			t.Errorf("Didn't get a value for [%s]", actual)
		}
		if !actual.Relative {
			t.Errorf("[%s] was not relative", tt.spath)
		}
		if actual.SegmentId != tt.seg_id {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.seg_id, actual.SegmentId)
		}
		if actual.IdValue != tt.qual {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.qual, actual.IdValue)
		}
		if actual.ElementIdx != tt.eleidx {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.eleidx, actual.ElementIdx)
		}
		if actual.SubelementIdx != tt.subeleidx {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.subeleidx, actual.SubelementIdx)
		}
		if strings.Join(actual.Loops, "/") != strings.Join(tt.loops, "/") {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", strings.Join(tt.loops, "/"), strings.Join(actual.Loops, "/"))
		}
		path := actual.String()
		if path != tt.spath {
			t.Errorf("Didn't get expected result [%s], instead got [%s]", tt.spath, path)
		}
	}
}

func stringSliceEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/*
class RefDes(unittest.TestCase):

    def testLoopOK1(self):
        path_str = "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/2000B/2300/2400"
        path = pyx12.path.X12Path(path_str)
        self.assertEqual(path_str, path.format())
        self.assertEqual(path.seg_id, None)
        self.assertEqual(path.loop_list[2], "ST_LOOP")

    def testLoopSegOK1(self):
        path_str = "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/2000B/2300/2400/SV2"
        path = pyx12.path.X12Path(path_str)
        self.assertEqual(path_str, path.format())
        self.assertEqual(path.seg_id, "SV2")



    def test_bad_rel_paths(self):
        bad_paths = [
            "AA/03",
            "BB/CC/03-2"
        ]
        for spath in bad_paths:
            self.assertRaises(pyx12.errors.X12PathError,
                              pyx12.path.X12Path, spath)

    def test_plain_loops(self):
        paths = [
            "ISA_LOOP/GS_LOOP",
            "GS_LOOP",
            "ST_LOOP/DETAIL/2000",
            "GS_LOOP/ST_LOOP/DETAIL/2000A",
            "DETAIL/2000A/2000B",
            "2000A/2000B/2300",
            "2000B/2300/2400",
            "ST_LOOP/HEADER",
            "ISA_LOOP/GS_LOOP/ST_LOOP/HEADER/1000A",
            "GS_LOOP/ST_LOOP/HEADER/1000B"
        ]
        for spath in paths:
            plist = spath.split("/")
            rd = pyx12.path.X12Path(spath)
            self.assertEqual(rd.loop_list, plist,
                             "%s: %s != %s" % (spath, rd.loop_list, plist))


class AbsolutePath(unittest.TestCase):
    def test_paths_with_refdes(self):
        tests = [
            ("/AAA/TST", "TST", None, None, None, ["AAA"]),
            ("/B1000/TST02", "TST", None, 2, None, ["B1000"]),
            ("/1000B/TST03-2", "TST", None, 3, 2, ["1000B"]),
            ("/1000A/1000B/TST[AA]02", "TST", "AA", 2, None, [
                "1000A", "1000B"]),
            ("/AA/BB/CC/TST[1B5]03-1", "TST", "1B5", 3, 1, ["AA", "BB", "CC"]),
            ("/DDD/E1000/N102", "N1", None, 2, None, ["DDD", "E1000"]),
            ("/E1000/D322/N102-5", "N1", None, 2, 5, ["E1000", "D322"]),
            ("/BB/CC/N1[AZR]02", "N1", "AZR", 2, None, ["BB", "CC"]),
            ("/BB/CC/N1[372]02-5", "N1", "372", 2, 5, ["BB", "CC"])
        ]
        for (spath, seg_id, qual, eleidx, subeleidx, plist) in tests:
            rd = pyx12.path.X12Path(spath)
            self.assertEqual(rd.relative, False,
                             "%s: %s != %s" % (spath, rd.relative, False))
            self.assertEqual(rd.seg_id, seg_id,
                             "%s: %s != %s" % (spath, rd.seg_id, seg_id))
            self.assertEqual(rd.id_val, qual, "%s: %s != %s" %
                             (spath, rd.id_val, qual))
            self.assertEqual(rd.ele_idx, eleidx,
                             "%s: %s != %s" % (spath, rd.ele_idx, eleidx))
            self.assertEqual(rd.subele_idx, subeleidx, "%s: %s != %s" %
                             (spath, rd.subele_idx, subeleidx))
            self.assertEqual(rd.format(), spath,
                             "%s: %s != %s" % (spath, rd.format(), spath))
            self.assertEqual(rd.loop_list, plist,
                             "%s: %s != %s" % (spath, rd.loop_list, plist))

    def test_bad_paths(self):
        bad_paths = [
            "/AA/03",
            "/BB/CC/03-2"
        ]
        for spath in bad_paths:
            self.assertRaises(pyx12.errors.X12PathError,
                              pyx12.path.X12Path, spath)

    def test_plain_loops(self):
        paths = [
            "/ISA_LOOP/GS_LOOP",
            "/ISA_LOOP/GS_LOOP",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/2000B",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/2000B/2300",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/2000B/2300/2400",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/HEADER",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/HEADER/1000A",
            "/ISA_LOOP/GS_LOOP/ST_LOOP/HEADER/1000B"
        ]
        for spath in paths:
            plist = spath.split("/")[1:]
            rd = pyx12.path.X12Path(spath)
            self.assertEqual(rd.loop_list, plist,
                             "%s: %s != %s" % (spath, rd.loop_list, plist))


class Equality(unittest.TestCase):
    def test_equal1(self):
        p1 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        p2 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        self.assertEqual(p1, p2)
        self.assertEqual(p1.format(), p2.format())
        self.assertEqual(p1.__hash__(), p2.__hash__())

    def test_equal2(self):
        p1 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        p2 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/")
        p2.loop_list.append("2000A")
        self.assertEqual(p1, p2)
        self.assertEqual(p1.format(), p2.format())
        self.assertEqual(p1.__hash__(), p2.__hash__())

    def test_equal3(self):
        p1 = pyx12.path.X12Path("/AA/BB/CC/TST[1B5]03-1")
        p2 = pyx12.path.X12Path("/AA/BB/CC/AAA[1B5]03-1")
        p2.seg_id = "TST"
        self.assertEqual(p1, p2)
        self.assertEqual(p1.format(), p2.format())
        self.assertEqual(p1.__hash__(), p2.__hash__())

    def test_equal4(self):
        p1 = pyx12.path.X12Path("1000B/TST03-2")
        p2 = pyx12.path.X12Path("1000B/TST04-2")
        p2.ele_idx = 3
        self.assertEqual(p1, p2)
        self.assertEqual(p1.format(), p2.format())
        self.assertEqual(p1.__hash__(), p2.__hash__())


class NonEquality(unittest.TestCase):
    def test_nequal1(self):
        p1 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        p2 = pyx12.path.X12Path("ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        self.assertNotEqual(p1, p2)
        self.assertNotEqual(p1.format(), p2.format())
        self.assertNotEqual(p1.__hash__(), p2.__hash__())

    def test_nequal2(self):
        p1 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A")
        p2 = pyx12.path.X12Path("/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/")
        self.assertNotEqual(p1, p2)
        self.assertNotEqual(p1.format(), p2.format())
        self.assertNotEqual(p1.__hash__(), p2.__hash__())

    def test_nequal3(self):
        p1 = pyx12.path.X12Path("/AA/BB/CC/TST[1B5]03-1")
        p2 = pyx12.path.X12Path("/AA/BB/CC/AAA[1B5]03-1")
        self.assertNotEqual(p1, p2)
        self.assertNotEqual(p1.format(), p2.format())
        self.assertNotEqual(p1.__hash__(), p2.__hash__())

    def test_nequal4(self):
        p1 = pyx12.path.X12Path("1000B/TST03-2")
        p2 = pyx12.path.X12Path("1000B/TST04-2")
        self.assertNotEqual(p1, p2)
        self.assertNotEqual(p1.format(), p2.format())
        self.assertNotEqual(p1.__hash__(), p2.__hash__())


class Empty(unittest.TestCase):
    def test_not_empty_1(self):
        p1 = "1000B/TST03-2"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_not_empty_2(self):
        p1 = "/AA/BB/CC/AAA[1B5]03"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_not_empty_3(self):
        p1 = "GS_LOOP/ST_LOOP/DETAIL/2000A"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_not_empty_4(self):
        p1 = "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_not_empty_5(self):
        p1 = "/"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_not_empty_6(self):
        p1 = "/ISA_LOOP/GS_LOOP/ST_LOOP/DETAIL/2000A/"
        self.assertFalse(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is not empty" % (p1))

    def test_empty_1(self):
        p1 = ""
        a = pyx12.path.X12Path(p1)
        self.assertTrue(pyx12.path.X12Path(
            p1).empty(), "Path "%s" is empty" % (p1))
*/