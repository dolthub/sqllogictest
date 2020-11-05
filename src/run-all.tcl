#!/usr/bin/tclsh
# Copyright 2019-2020 Dolthub, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Run this script in the "src" subdirectory of sqllogictest, after first
# compiling the ./sqllogictest binary, in order to verify correct output
# of all historical test cases.
#

set starttime [clock seconds]

if {$tcl_platform(platform)=="unix"} {
  set BIN ./sqllogictest
} else {
  set BIN ./sqllogictest.exe
}
if {![file exec $BIN]} {
  error "$BIN does not exist or is not executable.  Run make."
}

# add all test case file in the $subdir subdirectory to the
# set of all test case files in the global tcase() array.
#
proc search_for_test_cases {subdir} {
  foreach nx [glob -nocomplain $subdir/*] {
    if {[file isdir $nx]} {
      search_for_test_cases $nx
    } elseif {[string match *.test $nx]} {
      set ::tcase($nx) 1
    }
  }
}
search_for_test_cases ../test

# Run the tests
#
set totalerr 0
set totaltest 0
set totalrun 0
foreach tx [lsort [array names tcase]] {
  foreach opt {0 0xfff} {
    set opt "integrity_check;optimizer=[expr {$opt+0}]"
    catch {
      exec $BIN -verify -parameter $opt $tx
    } res
    puts $res
    if {[regexp {(\d+) errors out of (\d+) tests} $res all nerr ntst]} {
      incr totalerr $nerr
      incr totaltest $ntst
    } else {
      error "test did not complete: $BIN -verify -parameter optimizer=$opt $tx"
    }
    incr totalrun
  }
}

set endtime [clock seconds]
set totaltime [expr {$endtime - $starttime}]
puts "$totalerr errors out of $totaltest tests and $totalrun invocations\
      in $totaltime seconds"
