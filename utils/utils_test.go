package utils

/*
 * AWS SSO CLI
 * Copyright (c) 2021-2022 Aaron Turner  <synfinatic at gmail dot com>
 *
 * This program is free software: you can redistribute it
 * and/or modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or with the authors permission any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsSuite(t *testing.T) {
	s := &UtilsTestSuite{}
	suite.Run(t, s)
}

func (suite *UtilsTestSuite) TestParseRoleARN() {
	t := suite.T()

	a, r, err := ParseRoleARN("arn:aws:iam::11111:role/Foo")
	assert.Equal(t, int64(11111), a)
	assert.Equal(t, "Foo", r)
	assert.NoError(t, err)

	_, _, err = ParseRoleARN("")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("arnFoo")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("arn:aws:iam::a:role/Foo")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("arn:aws:iam::000000011111:role")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("aws:iam:000000011111:role/Foo")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("invalid:arn:aws:iam::000000011111:role/Foo")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("arn:aws:iam::000000011111:role/Foo/Bar")
	assert.Error(t, err)

	_, _, err = ParseRoleARN("arn:aws:iam::-000000011111:role/Foo")
	assert.Error(t, err)
}

func willPanicMakeRoleARN() {
	MakeRoleARN(-1, "foo")
}

func (suite *UtilsTestSuite) TestMakeRoleARN() {
	t := suite.T()

	assert.Equal(t, "arn:aws:iam::000000011111:role/Foo", MakeRoleARN(11111, "Foo"))
	assert.Equal(t, "arn:aws:iam::000000711111:role/Foo", MakeRoleARN(711111, "Foo"))
	assert.Equal(t, "arn:aws:iam::000000000000:role/", MakeRoleARN(0, ""))

	assert.Panics(t, willPanicMakeRoleARN)
}

func willPanicMakeRoleARNs() {
	MakeRoleARNs("asdfasfdo", "foo")
}

func (suite *UtilsTestSuite) TestMakeRoleARNs() {
	t := suite.T()

	assert.Equal(t, "arn:aws:iam::000000011111:role/Foo", MakeRoleARNs("11111", "Foo"))
	assert.Equal(t, "arn:aws:iam::000000711111:role/Foo", MakeRoleARNs("711111", "Foo"))
	assert.Equal(t, "arn:aws:iam::000000711111:role/Foo", MakeRoleARNs("000711111", "Foo"))
	assert.Equal(t, "arn:aws:iam::000000000000:role/", MakeRoleARNs("0", ""))

	assert.Panics(t, willPanicMakeRoleARNs)
}

func (suite *UtilsTestSuite) TestEnsureDirExists() {
	t := suite.T()

	defer os.RemoveAll("./does_not_exist_dir")
	assert.NoError(t, EnsureDirExists("./testdata/role_tags.yaml"))
	assert.NoError(t, EnsureDirExists("./does_not_exist_dir/foo.yaml"))
	assert.NoError(t, EnsureDirExists("./does_not_exist_dir/bar/baz/foo.yaml"))

	assert.Error(t, EnsureDirExists("/foo/bar"))
}

func (suite *UtilsTestSuite) TestGetHomePath() {
	t := suite.T()

	assert.Equal(t, "/", GetHomePath("/"))
	assert.Equal(t, ".", GetHomePath("."))
	assert.Equal(t, "/foo/bar", GetHomePath("/foo/bar"))
	assert.Equal(t, "/foo/bar", GetHomePath("/foo////bar"))
	assert.Equal(t, "/bar", GetHomePath("/foo/../bar"))
	home, _ := os.UserHomeDir()
	x := filepath.Join(home, "foo/bar")
	assert.Equal(t, x, GetHomePath("~/foo/bar"))
}

func (suite *UtilsTestSuite) TestAccountToString() {
	t := suite.T()

	a, err := AccountIdToString(0)
	assert.NoError(t, err)
	assert.Equal(t, "000000000000", a)

	a, err = AccountIdToString(11111)
	assert.NoError(t, err)
	assert.Equal(t, "000000011111", a)

	a, err = AccountIdToString(999999999999)
	assert.NoError(t, err)
	assert.Equal(t, "999999999999", a)

	_, err = AccountIdToString(-1)
	assert.Error(t, err)

	_, err = AccountIdToString(-19999)
	assert.Error(t, err)
}

func (suite *UtilsTestSuite) TestAccountToInt64() {
	t := suite.T()

	_, err := AccountIdToInt64("")
	assert.Error(t, err)

	a, err := AccountIdToInt64("12345")
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), a)

	a, err = AccountIdToInt64("0012345")
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), a)

	_, err = AccountIdToInt64("0012345678912123344455323423423423424")
	assert.Error(t, err)

	_, err = AccountIdToInt64("abdcefgi")
	assert.Error(t, err)

	_, err = AccountIdToInt64("-1")
	assert.Error(t, err)
}

var checkValue string
var checkBrowser string

func testUrlOpener(url string) error {
	checkBrowser = "default browser"
	checkValue = url
	return nil
}

func testUrlOpenerWith(url, browser string) error {
	checkBrowser = browser
	checkValue = url
	return nil
}

func testClipboardWriter(url string) error {
	checkValue = url
	return nil
}

func testUrlOpenerError(url string) error {
	return fmt.Errorf("there was an error")
}

func testUrlOpenerWithError(url, browser string) error {
	return fmt.Errorf("there was an error")
}

func (suite *UtilsTestSuite) TestHandleUrl() {
	t := suite.T()

	assert.Error(t, HandleUrl("foo", "browser", "bar", "pre", "post"))

	// override the print method
	printWriter = new(bytes.Buffer)
	assert.NoError(t, HandleUrl("print", "browser", "bar", "pre", "post"))
	assert.Equal(t, "prebarpost", printWriter.(*bytes.Buffer).String())

	urlOpener = testUrlOpener
	urlOpenerWith = testUrlOpenerWith
	clipboardWriter = testClipboardWriter

	assert.NoError(t, HandleUrl("clip", "browser", "url", "pre", "post"))
	assert.Equal(t, "url", checkValue)

	assert.NoError(t, HandleUrl("open", "other-browser", "other-url", "pre", "post"))
	assert.Equal(t, "other-browser", checkBrowser)
	assert.Equal(t, "other-url", checkValue)

	assert.NoError(t, HandleUrl("open", "", "some-url", "pre", "post"))
	assert.Equal(t, "default browser", checkBrowser)
	assert.Equal(t, "some-url", checkValue)

	urlOpener = testUrlOpenerError
	assert.Error(t, HandleUrl("open", "", "url", "pre", "post"))

	urlOpenerWith = testUrlOpenerWithError
	assert.Error(t, HandleUrl("open", "foo", "url", "pre", "post"))

	clipboardWriter = testUrlOpenerError
	assert.Error(t, HandleUrl("clip", "", "url", "pre", "post"))
}

func (suite *UtilsTestSuite) TestParseTimeString() {
	t := suite.T()

	x, e := ParseTimeString("1970-01-01 00:00:00 +0000 GMT")
	assert.NoError(t, e)
	assert.Equal(t, int64(0), x)
}

func (suite *UtilsTestSuite) TestTimeRemain() {
	t := suite.T()

	x, e := TimeRemain(0, false)
	assert.NoError(t, e)
	assert.Equal(t, "Expired", x)

	d, _ := time.ParseDuration("5m")
	future := time.Now().Add(d)
	x, e = TimeRemain(future.Unix(), true)
	assert.NoError(t, e)
	assert.Equal(t, "   5m", x)

	x, e = TimeRemain(future.Unix(), false)
	assert.NoError(t, e)
	assert.Equal(t, "5m", x)

	d, _ = time.ParseDuration("5h5m")
	future = time.Now().Add(d)
	x, e = TimeRemain(future.Unix(), true)
	assert.NoError(t, e)
	assert.Equal(t, "5h 5m", x)

	x, e = TimeRemain(future.Unix(), false)
	assert.NoError(t, e)
	assert.Equal(t, "5h5m", x)
}
