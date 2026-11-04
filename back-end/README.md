```markdown
# SocialDownloader — Backend

## Overview

This directory is intended to contain the SocialDownloader backend. The source code and dependency manifests were not included with the supplied project content, so the backend framework, architecture, API routes, and exact startup command cannot be confirmed here. Update this section after reviewing the backend source.

## Architecture

The architecture cannot be documented accurately without the backend source. Document the actual components here, including:

- The service entry point and framework
- API routes/controllers and their responsibilities
- Download or Instagram integration services
- Configuration and required environment variables
- Storage, database, and external services, if any

## Running the service

The exact run command depends on the backend's package/dependency manifest, which was not provided. From this directory, check for a manifest such as `package.json`, `pyproject.toml`, or `requirements.txt`; follow its scripts or setup instructions to install dependencies and start the service. Add the verified install and start commands here.

## Instagram cookie file

For private-account downloads, use this only with an account you own or are authorized to access, and follow the account first. If the application supports cookie-file authentication, export your own Instagram session cookies locally from your browser and save them as `insta_cookie.txt` in the location expected by the application. Confirm that location in the source/configuration before running the service.

Cookie files grant access to your account. Do not commit `insta_cookie.txt`, share it, or include it in logs or screenshots. Add it to `.gitignore`, keep it local, and delete or rotate the session if it is exposed. Respect Instagram's terms and the account owner's privacy.
```
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/download" -ContentType "application/json" -Body '{"url": "post_url"}' => for download

Invoke-RestMethod -Method POST -Uri "http://localhost:8080/extract" -ContentType "application/json" -Body '{"url": "post_url"}' => for extract the post info
