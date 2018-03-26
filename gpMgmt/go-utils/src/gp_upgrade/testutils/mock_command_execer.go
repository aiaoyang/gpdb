package testutils

import (
	"strings"
	"sync"

	"gp_upgrade/helpers"
)

type FakeCommandExecer struct {
	command string
	args    []string
	calls   []string

	mu      sync.Mutex
	output  helpers.Command
	trigger chan struct{}
}

func (c *FakeCommandExecer) SetOutput(command helpers.Command) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.output = command
}

func (c *FakeCommandExecer) SetTrigger(trigger chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.trigger = trigger
}

func (c *FakeCommandExecer) Calls() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.calls
}

func (c *FakeCommandExecer) Command() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.command
}

func (c *FakeCommandExecer) Args() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.args
}

func (c *FakeCommandExecer) Exec(command string, args ...string) helpers.Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.command = command
	c.args = args

	calledWith := append([]string{command}, args...)
	c.calls = append(c.calls, strings.Join(calledWith, " "))

	if c.trigger != nil {
		<-c.trigger
	}

	return c.output
}

type FakeCommand struct {
	numInvocations int
	Err            error
	Out            []byte
	Trigger        chan struct{}
}

func (c *FakeCommand) Output() ([]byte, error) {
	c.numInvocations++

	return c.Out, c.Err
}

func (c *FakeCommand) CombinedOutput() ([]byte, error) {
	return c.Out, c.Err
}

func (c *FakeCommand) Start() error {
	return c.Err
}

func (c *FakeCommand) Run() error {
	if c.Trigger != nil {
		<-c.Trigger
	}

	return c.Err
}

func (c *FakeCommand) GetNumInvocations() int {
	return c.numInvocations
}