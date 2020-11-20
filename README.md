# kubecon-blog

My blog about my KubeCon talk, also my demo.

## Try it out

```bash
export ADMIN_PASSWORD="" # To access the add-post function
export GITHUB_TOKEN="" # Personal access token

faas-cli secret create admin-token --from-literal $ADMIN_PASSWORD
faas-cli secret create github-token --from-literal $GITHUB_TOKEN

faas-cli deploy
```

Then go ahead and configure your GitHub Actions as per main.yml, creating secrets for:

* OPENFAAS_GATEWAY
* OPENFAAS_GATEWAY_PASSWD
* DOCKER_USERNAME
* DOCKER_PASSWORD

Whenever you add a post by add-posts, it will trigger a git commit, which triggers the GitHub workflow, to publish a new version of the blog function.

* [add-post](add-post) - checkout the handler.go and see how it does a commit from the markdown you enter.
* [blog](blog) - the content for the Hugo blog
* [openfaas-hugo-template](https://github.com/utsavanand2/openfaas-hugo-template) - the template for use with faas-cli


## Additional setup instructions from scratch

You would need to change the file `kubecon-multi-arch-blog.yml` to point to the gateway of your OpenFaaS deployment.

General overview of the steps to make this run locally which can be mapped to Github Actions:

```bash
# Pull templates
faas-cli template pull https://github.com/utsavanand2/openfaas-hugo-template

# Create new function (not required in github actions)
faas-cli new kubecon-multi-arch-blog --lang hugo --prefix alexellis2

# cd into kubecon-multi-arch-blog (not required in github actions)
cd kubecon-multi-arch-blog

# Create a new site (not required in github actions)
hugo new site .

# IMPORTANT. Make sure to change the base URL in config.toml 
# from-> baseURL = "http://example.org/"
# to-> baseURL = "http://$IP:8080/function/kubecon-multi-arch-blog/"

# Pull the hugo theme (this cmd is run with a slight modification in github actions)
git submodule add https://github.com/budparr/gohugo-theme-ananke.git themes/ananke && \
echo 'theme = "ananke"' >> config.toml

# Create a new post (not required in github actions)
hugo new posts/my-blog-is-serverless-and-multi-arch.md

# cd back into the root of the project dir (not required in github actions)
cd ..

# Add the themes dir to .gitignore because it causes trouble in Github Actions
echo kubecon-multi-arch-blog/themes >> .gitignore

# Run shrinkwrap build
faas-cli build -f kubecon-multi-arch-blog.yml --shrinkwrap

# Build with Docker Buildx (this cmd is run with a slight modification in github actions)
docker buildx build --platform linux/amd64,linux/arm64 --output "type=image,push=true" -t alexellis2/kubecon-multi-arch-blog:latest --no-cache build/kubecon-multi-arch-blog/

# Deploy to OpenFaaS (this cmd is run with a slight modification in github actions)
faas-cli deploy -f kubecon-multi-arch-blog.yml

# Let's Thank Matias for his awesome blog post! Made our lives so much easy!
https://www.openfaas.com/blog/serverless-static-sites/
```
