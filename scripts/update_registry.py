import json
import sys
import os

def update_registry(checksums_file, registry_file, version):
    # Read checksums
    checksums = {}
    with open(checksums_file, 'r') as f:
        for line in f:
            parts = line.strip().split()
            if len(parts) == 2:
                checksums[parts[1]] = parts[0]

    # Read registry
    with open(registry_file, 'r') as f:
        registry = json.load(f)

    # Update registry
    updated = False
    for extension in registry.get('extensions', []):
        for ver in extension.get('versions', []):
            if ver['version'] == version:
                print(f"Updating version {version}")
                for platform, artifact in ver.get('artifacts', {}).items():
                    url = artifact['url']
                    filename = url.split('/')[-1]
                    if filename in checksums:
                        print(f"  Updating {platform} ({filename})")
                        artifact['checksum']['value'] = checksums[filename]
                        updated = True
                    else:
                        print(f"  Warning: Checksum not found for {filename}")

    if updated:
        with open(registry_file, 'w') as f:
            json.dump(registry, f, indent=2)
            f.write('\n') # Add trailing newline
        print("Registry updated successfully.")
    else:
        print("No updates made.")

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python update_registry.py <checksums_file> <registry_file> <version>")
        sys.exit(1)
    
    update_registry(sys.argv[1], sys.argv[2], sys.argv[3])
