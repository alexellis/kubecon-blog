# kubecon-blog
My blog about my KubeCon talk, also my demo

You would need to change the file `kubecon-multi-arch-blog.yml` to point to the gateway of your OpenFaaS deployment.


General overview of the steps to make this run locally which can be mapped to Github Actions:

```bash
# Pull templates
faas-cli template pull https://github.com/utsavanand2/openfaas-hugo-template

# Create new function (not required in github actions)
faas-cli new blog --lang hugo --prefix alexellis2

# cd into blog (not required in github actions)
cd blog

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
echo blog/themes >> .gitignore

# Run shrinkwrap build
faas-cli build -f blog.yml --shrinkwrap

# Build with Docker Buildx (this cmd is run with a slight modification in github actions)
docker buildx build --platform linux/amd64,linux/arm64 --output "type=image,push=true" -t alexellis2/blog:latest --no-cache build/blog/

# Deploy to OpenFaaS (this cmd is run with a slight modification in github actions)
faas-cli deploy -f blog.yml

# Let's Thank Matias for his awesome blog post! Made our lives so much easy!
https://www.openfaas.com/blog/serverless-static-sites/
```