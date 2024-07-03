document.addEventListener('DOMContentLoaded', function() {
    console.log('Popup loaded');

    function saveFormData() {
        const project = {
            name: document.getElementById('name').value,
            displayName: document.getElementById('displayName').value,
            description: document.getElementById('description').value,
            websites: document.getElementById('website').value ? [{ url: document.getElementById('website').value }] : [],
            github: document.getElementById('github').value ? [{ url: document.getElementById('github').value }] : [],
            social: {}
        };

        if (document.getElementById('twitter').value) {
            project.social.twitter = [{ url: document.getElementById('twitter').value }];
        }

        if (document.getElementById('telegram').value) {
            project.social.telegram = [{ url: document.getElementById('telegram').value }];
        }

        if (document.getElementById('mirror').value) {
            project.social.mirror = [{ url: document.getElementById('mirror').value }];
        }

        if (Object.keys(project.social).length === 0) {
            delete project.social;
        }

        chrome.storage.sync.set({ project: project }, function() {
            console.log('Project data saved:', project);
        });
    }

    chrome.storage.sync.get(['project'], function(result) {
        if (result.project) {
            document.getElementById('name').value = result.project.name;
            document.getElementById('displayName').value = result.project.displayName;
            document.getElementById('description').value = result.project.description;
            if (result.project.websites.length > 0) {
                document.getElementById('website').value = result.project.websites[0].url;
            }
            if (result.project.github.length > 0) {
                document.getElementById('github').value = result.project.github[0].url;
            }
            if (result.project.social) {
                if (result.project.social.twitter && result.project.social.twitter.length > 0) {
                    document.getElementById('twitter').value = result.project.social.twitter[0].url;
                }
                if (result.project.social.telegram && result.project.social.telegram.length > 0) {
                    document.getElementById('telegram').value = result.project.social.telegram[0].url;
                }
                if (result.project.social.mirror && result.project.social.mirror.length > 0) {
                    document.getElementById('mirror').value = result.project.social.mirror[0].url;
                }
            }
            console.log('Project data loaded:', result.project);
        }
    });

    document.querySelectorAll('input[type="text"]').forEach(input => {
        input.addEventListener('input', saveFormData);
    });

    document.getElementById('projectForm').addEventListener('submit', function(event) {
        event.preventDefault();

        const project = {
            name: document.getElementById('name').value,
            displayName: document.getElementById('displayName').value,
            description: document.getElementById('description').value,
            websites: document.getElementById('website').value ? [{ url: document.getElementById('website').value }] : [],
            github: document.getElementById('github').value ? [{ url: document.getElementById('github').value }] : [],
            social: {}
        };

        if (document.getElementById('twitter').value) {
            project.social.twitter = [{ url: document.getElementById('twitter').value }];
        }

        if (document.getElementById('telegram').value) {
            project.social.telegram = [{ url: document.getElementById('telegram').value }];
        }

        if (document.getElementById('mirror').value) {
            project.social.mirror = [{ url: document.getElementById('mirror').value }];
        }

        if (Object.keys(project.social).length === 0) {
            delete project.social;
        }

        console.log('Sending request to create project with data:', project);

        fetch('http://localhost:8080/createProject', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(project)
        })
        .then(response => response.json())
        .then(data => {
            console.log('Project created:', data);
            if (data.error) {
                alert('Error creating project: ' + data.error);
            } else {
                alert('Project created successfully: ' + data.message);

                // Enable the git buttons
                document.getElementById('gitAddBtn').disabled = false;
                document.getElementById('gitCommitBtn').disabled = false;
                document.getElementById('gitPullBtn').disabled = false;
                document.getElementById('gitPushBtn').disabled = false;
            }
        })
        .catch((error) => {
            console.error('Error creating project:', error);
            alert('Error creating project: ' + error.message);
        });
    });

    document.getElementById('gitAddBtn').addEventListener('click', function() {
        runGitCommand('git add --all');
    });

    document.getElementById('gitCommitBtn').addEventListener('click', function() {
        const commitMessage = prompt("Enter commit message:", "Add new project");
        if (commitMessage) {
            runGitCommand(`git commit -m "${commitMessage}"`);
        }
    });

    document.getElementById('gitPullBtn').addEventListener('click', function() {
        runGitCommand('git pull origin main --rebase');
    });

    document.getElementById('gitPushBtn').addEventListener('click', function() {
        runGitCommand('git push origin main');
    });

    function runGitCommand(command) {
        console.log(`Executing git command: ${command}`);
        fetch('http://localhost:8080/runGitCommand', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ cmd: command })
        })
        .then(response => response.json())
        .then(data => {
            console.log('Git command output:', data);
            if (data.error) {
                alert('Error executing git command: ' + data.error);
            } else {
                alert('Git command executed successfully: ' + data.message);
            }
        })
        .catch((error) => {
            console.error('Error executing git command:', error);
            alert('Error executing git command: ' + error.message);
        });
    }
});
