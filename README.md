# Prototype Summary: Standalone CLI for Minder Rule Testing

### Project Overview
For my LFX mentorship proposal, I built a functional prototype of a standalone command-line interface (CLI) to evaluate and test Minder rules locally. The goal of this prototype was to solve a core developer experience (DX) problem: allowing rule authors to mathematically verify their security policies against mocked data instantly, without needing to spin up a live Minder server or wait for a CI pipeline to execute.

### The Problem I Addressed
Currently, testing Minder rules requires either deploying them to a live environment or writing tests tightly coupled to Go's internal `testing` package (`go test`), which outputs raw logs. Furthermore, the legacy test schema requires developers to write stringified JSON inside YAML blocks (e.g., `http.body: '{"foo":"bar"}'`), which is prone to syntax errors and difficult to read. 

### What I Built
I engineered a `pytest`-style CLI wrapper (`minder test -f rule.yaml -t test.yaml`) that provides an offline, instant, and highly readable testing framework. 

#### Key architectural achievements of the prototype include:
1. **Direct `rtengine` Integration:** Instead of building a custom Rego or JQ parser, I successfully decoupled Minder’s actual internal `rtengine` and `tkv1` (TestKit) packages from the `go test` dependency. By embedding the upstream engine directly into my CLI, the prototype guarantees 100% evaluation accuracy. If a rule passes in the CLI, it is guaranteed to pass in production.
2. **Modern Test Schema with Backward Compatibility:** I designed a cleaner test schema using a `mock_ingest` block, allowing authors to write native YAML instead of stringified JSON. Crucially, I built a backward-compatibility layer that automatically intercepts and translates legacy test formats, ensuring existing tests in the `minder-rules-and-profiles` repository will run perfectly without modification.
3. **Multi-Ingestion Support:** The prototype seamlessly handles both `type: rest` rules (mocking HTTP/JSON API responses) and `type: git` rules (mounting local `.testdata` folders to mock physical file systems and workflows).
4. **Profile Simulation:** The tool correctly maps `def` blocks from the test cases to simulate user profiles, allowing for dynamic JQ evaluation against mocked environments.

### How I Built It & Resources Used
The prototype was built using **Go** and the **Cobra** CLI framework. To ensure real-world validity, I directly utilized resources from the Stacklok/Minder ecosystem:
* **Core Dependencies:** Imported `github.com/mindersec/minder@latest` to leverage the official `minderv1` Protobuf definitions, `rtengine`, and `tkv1`.
* **Rule Extraction:** I pulled real, production-grade rules directly from the `minder-rules-and-profiles` repository to drive my test-driven development.
* **Test Cases:** I created custom test suites for highly complex rules to prove the framework's capability, specifically targeting:
    * `branch_protection_allow_deletions.yaml` (Testing JQ logic, HTTP 404 fallbacks, and parameter injection).
    * `github_actions_allowed.yaml` (Testing Profile-driven evaluations).
    * `grype_github_action_scan_container_image.yaml` (Testing Git filesystem ingestion using Rego).

### Images of test

**When test Passes**
<img width="1649" height="281" alt="swappy-20260409_182800" src="https://github.com/user-attachments/assets/f4cbbad8-2b95-4e0e-ad0d-4d08dac85f21" />

**When test fails**
<img width="1649" height="262" alt="swappy-20260409_182800 (Edited)" src="https://github.com/user-attachments/assets/a47cf13b-2196-4f10-a687-e629ce02ae67" />

