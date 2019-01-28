// Copyright 2018 DREP Foundation Ltd.
// This file is part of the drep-cli library.
//
// The drep-cli library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The drep-cli library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the drep-cli library. If not, see <http://www.gnu.org/licenses/>.

// Package deps contains the console JavaScript dependencies Go embedded.
package deps

// eth version v0.20.6
//go:generate go-bindata -nometadata -pkg deps -o bindata.go bignumber.js drep.js
//go:generate gofmt -w -s bindata.go
