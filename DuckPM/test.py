import click
import requests
import json
import re

# Base URL of the Quacko Retro Hub API
BASE_URL = "https://quacko-retro-hub-api.vercel.app"

# UUID validation regex
UUID_REGEX = r'^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'

# Helper function to validate UUID
def is_valid_uuid(value):
    return bool(re.match(UUID_REGEX, value, re.IGNORECASE))

# Helper function to handle API responses
def handle_response(response):
    try:
        response.raise_for_status()
        return response.json()
    except requests.exceptions.HTTPError as e:
        try:
            error_message = response.json().get('message', str(e))
        except ValueError:
            error_message = str(e)
        click.echo(f"Error: {error_message}")
        return None
    except requests.exceptions.RequestException as e:
        click.echo(f"Request failed: {str(e)}")
        return None

@click.group()
def cli():
    """Quacko Retro Hub CLI - Interact with vintage libraries via the API"""
    pass

@cli.command()
@click.option('--email', prompt='Email', help='Your email')
@click.option('--password', prompt='Password', hide_input=True, help='Your password')
def register(email, password):
    """Register a new user on Quacko Retro Hub"""
    payload = {"email": email, "password": password}
    response = requests.post(f"{BASE_URL}/register", json=payload)
    result = handle_response(response)
    if result:
        click.echo(f"✅ {result['message']}")
        click.echo(json.dumps(result.get('data', {}), indent=2))

@cli.command()
@click.option('--email', prompt='Email', help='Your email')
@click.option('--password', prompt='Password', hide_input=True, help='Your password')
def login(email, password):
    """Login to your Quacko Retro Hub account"""
    payload = {"email": email, "password": password}
    response = requests.post(f"{BASE_URL}/login", json=payload)
    result = handle_response(response)
    if result:
        click.echo(f"✅ {result['message']}")
        click.echo(json.dumps(result.get('data', {}), indent=2))

@cli.command()
def list_repos():
    """List all retro libraries on Quacko Retro Hub"""
    response = requests.get(f"{BASE_URL}/repos")
    result = handle_response(response)
    if result:
        click.echo("📚 Available Libraries:")
        click.echo(json.dumps(result.get('data', []), indent=2))

@cli.command()
@click.argument('repo_id')
def get_repo(repo_id):
    """Get detailed info of a specific library by ID"""
    if not is_valid_uuid(repo_id):
        click.echo("❌ Invalid repository ID. Must be a UUID.")
        return
    response = requests.get(f"{BASE_URL}/repos/{repo_id}")
    result = handle_response(response)
    if result:
        click.echo("📦 Repository Details:")
        click.echo(json.dumps(result.get('data', {}), indent=2))

@cli.command()
@click.option('--username', prompt='Email (username)', help='Your email')
@click.option('--password', prompt='Password', hide_input=True, help='Your password')
@click.option('--name', prompt='Name', help='Library name')
@click.option('--description', prompt='Description', help='Short description')
@click.option('--github-repo-link', prompt='GitHub Repo Link', help='Link to GitHub repo')
@click.option('--license', prompt='License', help='License (e.g., MIT)')
@click.option('--tags', prompt='Tags', help='Comma-separated tags')
@click.option('--version', prompt='Version', help='Version number (e.g., 1.0.0)')
def create_repo(username, password, name, description, github_repo_link, license, tags, version):
    """Publish a new retro library to Quacko Retro Hub"""
    payload = {
        "name": name,
        "description": description,
        "github_repo_link": github_repo_link,
        "license": license,
        "tags": [tag.strip() for tag in tags.split(',')],
        "version": version,
        "username": username,
        "password": password
    }
    response = requests.post(f"{BASE_URL}/repos", json=payload)
    result = handle_response(response)
    if result:
        click.echo("🚀 Library created successfully!")
        click.echo(json.dumps(result.get('data', {}), indent=2))

@cli.command()
@click.argument('repo_id')
@click.option('--username', prompt='Email (username)', help='Your email')
@click.option('--password', prompt='Password', hide_input=True, help='Your password')
def delete_repo(repo_id, username, password):
    """Delete a library (only if you're the author)"""
    if not is_valid_uuid(repo_id):
        click.echo("❌ Invalid repository ID. Must be a UUID.")
        return
    payload = {"username": username, "password": password}
    response = requests.delete(f"{BASE_URL}/repos/{repo_id}", json=payload)
    result = handle_response(response)
    if result:
        click.echo("🗑️ Repository deleted successfully!")

if __name__ == "__main__":
    cli()
