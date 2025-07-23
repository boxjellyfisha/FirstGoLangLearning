FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    lsb-release \
    golang-go \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Docker CLI
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

RUN echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
    $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update && apt-get install -y \
    docker-ce-cli \
    && rm -rf /var/lib/apt/lists/*

# Install Ansible and Docker collection
RUN pip install ansible

# Install community.docker collection
RUN ansible-galaxy collection install community.docker

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /ansible

# After the container is created, run the ansible-playbook command
# you can dynamically add arguments to the command
# example: docker run -it --rm src-ansible ansible-playbook -i ansible/inventory.ini ansible/do.yml -e "environment=prod"
ENTRYPOINT ["ansible-playbook"]
