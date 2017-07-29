package cprofile

import (
	"fmt"

	"github.com/google/pprof/profile"
)

// Process represents a profiled process
type Process struct {
	timeline *Timeline
}

// Timeline represents a goroutine
type Timeline struct {
	start int
	calls []Call
}

type Call struct {
	location int64
}

// NewProcess takes a pprof profile and translates it into a process
func NewProcess(profile *profile.Profile) *Process {
	p := &Process{}

	// // Create low levelest level calls.  These will be from main or go routines.
	// for _, s := range profile.Sample {

	// }

	fmt.Printf("-------\n\nMappings: %d\n", len(profile.Mapping))
	for _, m := range profile.Mapping {
		fmt.Printf("%d: File: %s / '%s' / %d / %d / %d\n", m.ID, m.File, m.BuildID, m.Limit, m.Offset, m.Start)
	}

	fmt.Printf("Comments: (%d)\n", len(profile.Comments))
	for _, c := range profile.Comments {
		fmt.Printf("%s\n", c)
	}

	fmt.Printf("\nSamples: %d\n", len(profile.Sample))
	for _, s := range profile.Sample {
		// Get first one
		l := s.Location[len(s.Location)-1]
		fmt.Printf("First loc: 0x%08x / ID: %d / mapping: %d\n", l.Address, l.ID, l.Mapping.ID)

		fmt.Printf("Labels: (%d)\n", len(s.Label))
		for k, ls := range s.Label {
			fmt.Printf("%s...\n", k)
			for k, l := range ls {
				fmt.Printf("\t%d: %s\n", k, l)
			}
		}

		fmt.Printf("Num Labels: (%d)\n", len(s.NumLabel))
		for k, v := range s.NumLabel {
			fmt.Printf("%s...\n", k)
			for k, v := range v {
				fmt.Printf("\t%d: %d\n", k, v)
			}
		}

		fmt.Printf("Value: (%d)\n", len(s.Value))
		for k, v := range s.Value {
			fmt.Printf("%d: %d\n", k, v)
		}

		fmt.Printf("\n")
	}
	return p
}
