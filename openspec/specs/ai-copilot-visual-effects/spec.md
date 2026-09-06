## Purpose

Defines the dynamic visual effects, glowing border beams, chromatic metal liquid surfaces, and thinking-orb states for the 0things AI Copilot frontend interface based on libraries.dev.

## Requirements

### Requirement: Chromatic Metal Floating Trigger Button
The floating AI Copilot trigger button SHALL be wrapped with `MetalFx` utilizing the `chromatic` preset to render real-time WebGL liquid metal reflections and interactive illumination halo, with theme responsiveness for dark and light modes.

#### Scenario: Rendering chromatic metal button
- **WHEN** the user views the AI Copilot floating trigger button on any authenticated page
- **THEN** the button displays a chromatic fluid metallic halo around its border, adjusting shader palette dynamically based on current theme (dark/light)

### Requirement: Full-Spectrum Animated Border Beam for Assistant Modal
The assistant modal dialog window SHALL render a continuous animated glowing border beam using `BorderBeam` with `colorful` color variant and medium size (`size="md"`).

#### Scenario: Opening assistant modal with border beam
- **WHEN** the user clicks the floating trigger button to expand the assistant modal
- **THEN** the modal window border animates with a full-spectrum flowing rainbow beam around its rounded container

### Requirement: Interactive Thought-Orb Loading Indicators
The AI Copilot interface SHALL integrate `ThinkingOrb` components to reflect granular lifecycle states:
1. Welcome screen hero displays a 64px thought orb in `breathing` state.
2. Active reasoning/thought blocks display a 20px thought orb in `solving` state while running, transitioning to `breathing` when complete.
3. IoT Tool execution blocks display a 20px thought orb in `searching` state while executing queries.
4. Assistant stream cursor indicator displays a subtle thought orb.

#### Scenario: Reasoning during model generation
- **WHEN** the AI model is streaming chain-of-thought reasoning tokens
- **THEN** the reasoning trigger displays a 20px `ThinkingOrb` in `solving` state with animated rotation

#### Scenario: IoT tool query execution
- **WHEN** the assistant invokes an IoT device or telemetry query tool
- **THEN** the tool group / fallback trigger renders a 20px `ThinkingOrb` in `searching` state

