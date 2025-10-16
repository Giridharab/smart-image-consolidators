# Smart Image Consolidator

This tool automates Dockerfile analysis in GitHub PRs. It provides:

1. **Canonical base image suggestions**  
2. **Real-time CPU/Memory/Storage/Cost metrics**  
3. **Security scan with AnchoreCTL**  
4. **PR comments with full results**

When you open a PR with these Dockerfiles, the Smart Image Consolidator workflow will:
* Detect the Dockerfiles
* Build temporary containers
* Measure CPU/Memory/Storage
* Suggest canonical base images (python:3.11-slim → artifactory.devhub-cloud.cisco.com/python:3.11-slim)
* Run AnchoreCTL for security scan
* Post a full comment on the PR


---

## **Setup**

1. Clone repo:

```bash
git clone https://github.com/Giridharab/smart-image-consolidator.git
cd smart-image-consolidator

2.Install Go 1.24+ and Docker CLI v20.x (Ubuntu runners already have Docker).
3.Create a GitHub personal access token (PAT) with repo permissions and store it in your repo secrets as GH_PAT.

How it works
1.The GitHub Actions workflow triggers on PR open/update.
2.Go scans all Dockerfiles in the PR.
3.For each Dockerfile:
 * Builds a temporary container to measure real-time CPU/Memory.
 * Inspects image for storage size.
 * Suggests canonical base image from canonical_bases.yaml.
 * Runs AnchoreCTL security scan.
4.Posts a single comment in the PR with all results.

Running locally

```
CanonicalBases:
  - Original: "python:3.11-slim"
    Suggested: "artifactory.devhub-cloud.cisco.com/python:3.11-slim"
```
* Ensure Docker daemon is running locally.
* AnchoreCTL must be installed (anchore --version).

Adding Canonical Base Images

Edit configs/canonical_bases.yaml:
```
CanonicalBases:
  - Original: "python:3.11-slim"
    Suggested: "artifactory.devhub-cloud.cisco.com/python:3.11-slim"
```

Workflow in Action
* The workflow is defined in .github/workflows/smart-image-consolidator.yml.
* Triggers automatically on pull requests:
```
on:
  pull_request:
    types: [opened, synchronize, reopened]
```
* Checks out the code, sets PR variables, builds Go binary, and runs consolidator.

Output Example
```
### Dockerfile: Dockerfile.python

CPU: 120.00%
Memory: 150MiB
Storage: 50.12MiB
Estimated Cost: $0.01
Canonical Base Suggestion: artifactory.devhub-cloud.cisco.com/python:3.11-slim

#### Security Scan (AnchoreCTL)
Vulnerability Summary:
- CVE-2025-XXXX: High
- CVE-2025-YYYY: Medium

```

Cleanup
* Temporary containers are removed automatically.
* Docker images can be pruned locally:
```
docker image prune -f
```



