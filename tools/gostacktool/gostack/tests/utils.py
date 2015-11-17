# Copyright 2015 The Serviced Authors.
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

import gostack


def populateGoroutine(goroutine, *args):
    for i in range(len(args)):
        if i == 0:
            warnings = goroutine.parseLine(args[0])
        elif (i % 2) == 1:
            stackframe = gostack.StackFrame()
            warnings += stackframe.parseFunctionLine(args[i])
        else:
            warnings += stackframe.parseFileLine(args[i])
            goroutine.addFrame(stackframe)
    return warnings