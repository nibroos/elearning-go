# Define the source of your .env and .env.example files
ENV_SOURCE=docker/.env
ENV_EXAMPLE_SOURCE=docker/.env.example

# Define the services you want to copy the .env files to
SERVICES=service users-service master-service activities-service

# Copy the .env and .env.example files to each service
copy-env:
	if [ ! -f $(ENV_SOURCE) ]; then \
		echo "File $(ENV_SOURCE) does not exist. Creating..."; \
		cp $(ENV_EXAMPLE_SOURCE) $(ENV_SOURCE); \
	fi
	
	@for service in $(SERVICES); do \
		if [ ! -d $$service ]; then \
			echo "Directory $$service does not exist. Creating..."; \
			mkdir -p $$service; \
		fi; \
		if [ ! -f $$service/.env ]; then \
			echo "Copying .env to $$service"; \
			cp $(ENV_SOURCE) $$service/.env; \
		fi; \
		if [ ! -f $$service/.env.example ]; then \
			echo "Copying .env.example to $$service"; \
			cp $(ENV_EXAMPLE_SOURCE) $$service/.env.example; \
		fi; \
	done
