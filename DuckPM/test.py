import requests
import json
import argparse
import re
import getpass

BASE_URL = "https://quacko-retro-hub-api.vercel.app"
UUID_REGEX = r'^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'

def is_valid_uuid(value):
    return bool(re.match(UUID_REGEX, value, re.IGNORECASE))

def handle_response(response):
    try:
        response.raise_for_status()
        return response.json()
    except requests.exceptions.HTTPError as e:
        try:
            error_message = response.json().get('message', str(e))
        except ValueError:
            error_message = str(e)
        print(f"[ERROR] {error_message}")
        return None
    except requests.exceptions.RequestException as e:
        print(f"[REQUEST FAILED] {str(e)}")
        return None

def register():
    email = input("Email: ")
    password = getpass.getpass("Password: ")
    payload = {"email": email, "password": password}
    response = requests.post(f"{BASE_URL}/register", json=payload)
    result = handle_response(response)
    if result:
        print("✅", result['message'])
        print(json.dumps(result.get('data', {}), indent=2))

def login():
    email = input("Email: ")
    password = getpass.getpass("Password: ")
    payload = {"email": email, "password": password}
    response = requests.post(f"{BASE_URL}/login", json=payload)
    result = handle_response(response)
    if result:
        print("✅", result['message'])
        print(json.dumps(result.get('data', {}), indent=2))

def list_repos():
    response = requests.get(f"{BASE_URL}/repos")
    result = handle_response(response)
    if result:
        print("📚 Available Libraries:")
        print(json.dumps(result.get('data', []), indent=2))

def get_repo(repo_id):
    if not is_valid_uuid(repo_id):
        print("❌ Invalid repository ID. Must be a UUID.")
        return
    response = requests.get(f"{BASE_URL}/repos/{repo_id}")
    result = handle_response(response)
    if result:
        print("📦 Repository Details:")
        print(json.dumps(result.get('data', {}), indent=2))

def create_repo():
    username = input("Email (username): ")
    password = getpass.getpass("Password: ")
    name = input("Name: ")
    description = input("Description: ")
    github_repo_link = input("GitHub Repo Link: ")
    license = input("License (e.g., MIT): ")
    tags = input("Tags (comma separated): ").split(',')
    version = input("Version (e.g., 1.0.0): ")
    
    payload = {
        "name": name,
        "description": description,
        "github_repo_link": github_repo_link,
        "license": license,
        "tags": [tag.strip() for tag in tags],
        "version": version,
        "username": username,
        "password": password
    }
    response = requests.post(f"{BASE_URL}/repos", json=payload)
    result = handle_response(response)
    if result:
        print("🚀 Library created successfully!")
        print(json.dumps(result.get('data', {}), indent=2))

def delete_repo(repo_id):
    if not is_valid_uuid(repo_id):
        print("❌ Invalid repository ID. Must be a UUID.")
        return
    username = input("Email (username): ")
    password = getpass.getpass("Password: ")
    payload = {"username": username, "password": password}
    response = requests.delete(f"{BASE_URL}/repos/{repo_id}", json=payload)
    result = handle_response(response)
    if result:
        print("🗑️ Repository deleted successfully!")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="🦆 Quacko Retro Hub CLI (No Click, No Bullshit)")
    subparsers = parser.add_subparsers(dest="command")

    subparsers.add_parser("register")
    subparsers.add_parser("login")
    subparsers.add_parser("list-repos")

    get_parser = subparsers.add_parser("get-repo")
    get_parser.add_argument("repo_id")

    subparsers.add_parser("create-repo")

    del_parser = subparsers.add_parser("delete-repo")
    del_parser.add_argument("repo_id")

    args = parser.parse_args()

    if args.command == "register":
        register()
    elif args.command == "login":
        login()
    elif args.command == "list-repos":
        list_repos()
    elif args.command == "get-repo":
        get_repo(args.repo_id)
    elif args.command == "create-repo":
        create_repo()
    elif args.command == "delete-repo":
        delete_repo(args.repo_id)
    else:
        parser.print_help()
