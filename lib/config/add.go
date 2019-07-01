// Copyright (c) 2019 Ashley Jeffs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package config

import (
	"errors"
	"fmt"

	"github.com/Jeffail/benthos/lib/cache"
	"github.com/Jeffail/benthos/lib/condition"
	"github.com/Jeffail/benthos/lib/config"
	"github.com/Jeffail/benthos/lib/input"
	"github.com/Jeffail/benthos/lib/output"
	"github.com/Jeffail/benthos/lib/processor"
	"github.com/Jeffail/benthos/lib/ratelimit"
)

//------------------------------------------------------------------------------

// AddInput inserts a default input of a type to an existing config.
func AddInput(cType string, conf *config.Type) error {
	if cType != "benthos_lab" {
		if _, ok := input.Constructors[cType]; !ok {
			return fmt.Errorf("input type '%v' not recognised", cType)
		}
	}
	inputConf := input.NewConfig()
	inputConf.Type = cType

	if conf.Input.Type != input.TypeBroker {
		currentInput := conf.Input
		brokerInput := input.NewConfig()
		brokerInput.Type = input.TypeBroker
		brokerInput.Broker.Inputs = append(brokerInput.Broker.Inputs, currentInput)
		conf.Input = brokerInput
		if cType == input.TypeBroker {
			return nil
		}
	}
	conf.Input.Broker.Inputs = append(conf.Input.Broker.Inputs, inputConf)
	return nil
}

// AddProcessor inserts a default processor of a type to an existing config.
func AddProcessor(cType string, conf *config.Type) error {
	if _, ok := processor.Constructors[cType]; !ok {
		return fmt.Errorf("processor type '%v' not recognised", cType)
	}
	procConf := processor.NewConfig()
	procConf.Type = cType

	conf.Pipeline.Processors = append(conf.Pipeline.Processors, procConf)
	return nil
}

// AddCondition inserts a filter_parts processor with a default condition of a
// type to an existing config.
func AddCondition(cType string, conf *config.Type) error {
	if _, ok := condition.Constructors[cType]; !ok {
		return fmt.Errorf("condition type '%v' not recognised", cType)
	}
	condConf := condition.NewConfig()
	condConf.Type = cType

	procConf := processor.NewConfig()
	procConf.Type = processor.TypeFilterParts
	procConf.FilterParts.Config = condConf

	conf.Pipeline.Processors = append(conf.Pipeline.Processors, procConf)
	return nil
}

// AddOutput inserts a default output of a type to an existing config.
func AddOutput(cType string, conf *config.Type) error {
	if cType != "benthos_lab" {
		if _, ok := output.Constructors[cType]; !ok {
			return fmt.Errorf("output type '%v' not recognised", cType)
		}
	}
	outputConf := output.NewConfig()
	outputConf.Type = cType

	if conf.Output.Type != output.TypeBroker {
		currentOutput := conf.Output
		brokerOutput := output.NewConfig()
		brokerOutput.Type = output.TypeBroker
		brokerOutput.Broker.Outputs = append(brokerOutput.Broker.Outputs, currentOutput)
		conf.Output = brokerOutput
		if cType == output.TypeBroker {
			return nil
		}
	}
	conf.Output.Broker.Outputs = append(conf.Output.Broker.Outputs, outputConf)
	return nil
}

// AddCache inserts a default cache of a type to an existing config.
func AddCache(cType string, conf *config.Type) error {
	if _, ok := cache.Constructors[cType]; !ok {
		return fmt.Errorf("cache type '%v' not recognised", cType)
	}
	cacheConf := cache.NewConfig()
	cacheConf.Type = cType

	var cacheID string
	for i := 0; i < 10000; i++ {
		var candidate string
		if i == 0 {
			candidate = "example"
		} else {
			candidate = fmt.Sprintf("example%v", i)
		}
		if _, exists := conf.Manager.Caches[candidate]; !exists {
			cacheID = candidate
			break
		}
	}
	if len(cacheID) == 0 {
		return errors.New("what the hell are you doing?")
	}

	conf.Manager.Caches[cacheID] = cacheConf
	return nil
}

// AddRatelimit inserts a default rate limit of a type to an existing config.
func AddRatelimit(cType string, conf *config.Type) error {
	if _, ok := ratelimit.Constructors[cType]; !ok {
		return fmt.Errorf("ratelimit type '%v' not recognised", cType)
	}
	ratelimitConf := ratelimit.NewConfig()
	ratelimitConf.Type = cType

	var ratelimitID string
	for i := 0; i < 10000; i++ {
		var candidate string
		if i == 0 {
			candidate = "example"
		} else {
			candidate = fmt.Sprintf("example%v", i)
		}
		if _, exists := conf.Manager.RateLimits[candidate]; !exists {
			ratelimitID = candidate
			break
		}
	}
	if len(ratelimitID) == 0 {
		return errors.New("what the hell are you doing?")
	}

	conf.Manager.RateLimits[ratelimitID] = ratelimitConf
	return nil
}

//------------------------------------------------------------------------------
