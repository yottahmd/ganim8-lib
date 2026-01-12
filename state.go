package ganim8

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// State represents the state of an animation.
type State struct {
	Reference *Animation
	position  int
	timer     time.Duration
	status    Status
}

func NewState(reference *Animation) State {
	return State{
		Reference: reference,
		position:  0,
		timer:     0,
		status:    Playing,
	}
}

// Update updates the animation state.
func (state *State) Update() {
	state.UpdateWithDelta(DefaultDelta)
}

// UpdateWithDelta updates the animation state with the specified delta.
func (state *State) UpdateWithDelta(elapsedTime time.Duration) {
	if state.status != Playing || state.Reference.sprite.length <= 1 {
		return
	}
	state.timer += elapsedTime
	loops := state.timer / state.Reference.totalDuration
	if loops != 0 {
		state.timer = state.timer - state.Reference.totalDuration*loops
		(state.Reference.onLoop)(state.Reference, int(loops))
	}
	state.position = seekFrameIndex(state.Reference.intervals, state.timer)
}

func (state *State) Reset() {
	state.position = 0
	state.timer = 0
}

func (state *State) IsEnd() bool {
	if state.status == Paused && state.position == state.Reference.sprite.length-1 {
		return true
	}
	return false
}

// Status returns the status of the animation state.
func (state *State) Status() Status {
	return state.status
}

// Pause pauses the animation state.
func (state *State) Pause() {
	state.status = Paused
}

// Position returns the current position of the frame.
// The position counts from 1 (not 0).
func (state *State) Position() int {
	return state.position + 1
}

// Timer returns the current accumulated times of the current frame.
func (state *State) Timer() time.Duration {
	return state.timer
}

// GoToFrame sets the position of the animation state and
// sets the timer at the start of the frame.
func (state *State) GoToFrame(position int) {
	state.position = position - 1
	state.timer = state.Reference.intervals[state.position]
}

// PauseAtEnd pauses the animation state and sets the position
// to the last frame.
func (state *State) PauseAtEnd() {
	state.position = state.Reference.sprite.length - 1
	state.timer = state.Reference.totalDuration
	state.Pause()
}

// PauseAtStart pauses the animation state and sets the position
// to the first frame.
func (state *State) PauseAtStart() {
	state.position = 0
	state.timer = 0
	state.status = Paused
}

// Resume resumes the animation state.
func (state *State) Resume() {
	state.status = Playing
}

// Draw draws the animation with the specified option parameters.
func (state *State) Draw(screen *ebiten.Image, opts *DrawOptions) {
	state.Reference.sprite.Draw(screen, state.position, opts)
}

// DrawWithShader draws the animation with the specified option parameters.
func (state *State) DrawWithShader(screen *ebiten.Image, opts *DrawOptions, shaderOpts *ShaderOptions) {
	state.Reference.sprite.DrawWithShader(screen, state.position, opts, shaderOpts)
}
