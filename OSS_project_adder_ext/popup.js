document.addEventListener('DOMContentLoaded', function() {
    console.log('Popup loaded');

    chrome.storage.local.get(['website', 'github', 'twitter', 'telegram', 'mirror'], function(result) {
        if (result.website) document.getElementById('website').value = result.website;
        if (result.github) document.getElementById('github').value = result.github;
        if (result.twitter) document.getElementById('twitter').value = result.twitter;
        if (result.telegram) document.getElementById('telegram').value = result.telegram;
        if (result.mirror) document.getElementById('mirror').value = result.mirror;
        if (result.discord) document.getElementById('discord').value = result.discord;
    });

    chrome.storage.local.get(['persistentFiles', 'completedFiles'], function(result) {
        if (result.persistentFiles) {
            const filesList = document.getElementById('addedFilesList');
            filesList.innerHTML = '';
            result.persistentFiles.forEach(file => {
                const li = document.createElement('li');
                li.style.display = 'flex';
                li.style.alignItems = 'center';
                li.style.marginBottom = '8px';
                
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.checked = result.completedFiles && result.completedFiles.includes(file);
                checkbox.style.marginRight = '8px';
                checkbox.addEventListener('change', (e) => {
                    chrome.storage.local.get(['completedFiles'], function(storage) {
                        const completed = new Set(storage.completedFiles || []);
                        if (e.target.checked) {
                            completed.add(file);
                        } else {
                            completed.delete(file);
                        }
                        chrome.storage.local.set({ 'completedFiles': Array.from(completed) });
                    });
                });
                
                const a = document.createElement('a');
                a.href = '#';
                a.textContent = file;
                if (result.completedFiles && result.completedFiles.includes(file)) {
                    a.style.textDecoration = 'line-through';
                    a.style.color = '#666';
                }
                a.addEventListener('click', (e) => {
                    e.preventDefault();
                    fetchFileContent(file);
                });
                
                li.appendChild(checkbox);
                li.appendChild(a);
                filesList.appendChild(li);
            });
        }
    });

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

        if (document.getElementById('discord').value) {
            project.social.discord = [{ url: document.getElementById('discord').value }];
        }

        if (Object.keys(project.social).length === 0) {
            delete project.social;
        }

        chrome.storage.sync.set({ project: project }, function() {
            console.log('Project data saved:', project);
        });
    }

    function loadFormData() {
        chrome.storage.sync.get(['project'], function(result) {
            if (result.project) {
                document.getElementById('name').value = result.project.name || '';
                document.getElementById('displayName').value = result.project.displayName || '';
                document.getElementById('description').value = result.project.description || '';
                if (result.project.websites && result.project.websites.length > 0) {
                    document.getElementById('website').value = result.project.websites[0].url || '';
                }
                if (result.project.github && result.project.github.length > 0) {
                    document.getElementById('github').value = result.project.github[0].url || '';
                }
                if (result.project.social) {
                    if (result.project.social.twitter && result.project.social.twitter.length > 0) {
                        document.getElementById('twitter').value = result.project.social.twitter[0].url || '';
                    }
                    if (result.project.social.telegram && result.project.social.telegram.length > 0) {
                        document.getElementById('telegram').value = result.project.social.telegram[0].url || '';
                    }
                    if (result.project.social.mirror && result.project.social.mirror.length > 0) {
                        document.getElementById('mirror').value = result.project.social.mirror[0].url || '';
                    }
                    if (result.project.social.discord && result.project.social.discord.length > 0) {
                        document.getElementById('discord').value = result.project.social.discord[0].url || '';
                    }
                }
                console.log('Project data loaded:', result.project);
            }
        });
        chrome.storage.local.get(['selectedUrl'], function(result) {
            if (result.selectedUrl) {
                fillUrlField(result.selectedUrl);
                // Clear the stored URL after using it
                chrome.storage.local.remove('selectedUrl');
            }
        });
    }

    function fillUrlField(url) {
        // Determine which field to fill based on the URL
        if (url.includes('github.com')) {
            document.getElementById('github').value = url;
        } else if (url.includes('twitter.com') || url.includes('x.com')) {
            document.getElementById('twitter').value = url;
        } else if (url.includes('t.me')) {
            document.getElementById('telegram').value = url;
        } else if (url.includes('mirror.xyz')) {
            document.getElementById('mirror').value = url;
        } else if (url.includes('discord.com') || url.includes('discord.gg')) {
            document.getElementById('discord').value = url;
        } else {
            document.getElementById('website').value = url;
        }
    }

    function updateAddedFilesList() {
        chrome.storage.local.get(['completedFiles'], function(storage) {
            const completedFiles = new Set(storage.completedFiles || []);
            
            fetch('http://localhost:8080/getAddedFiles')
            .then(response => response.json())
            .then(data => {
                const filesList = document.getElementById('addedFilesList');
                if (!filesList) {
                    const filesSection = document.getElementById('addedFilesSection');
                    const newList = document.createElement('ul');
                    newList.id = 'addedFilesList';
                    filesSection.appendChild(newList);
                }
                filesList.innerHTML = '';
                
                chrome.storage.local.set({ 'persistentFiles': data.files });
                
                data.files.forEach(file => {
                    const li = document.createElement('li');
                    li.style.cssText = `
                        display: flex;
                        align-items: center;
                        margin: 8px 0;
                        padding: 8px;
                        border-radius: 4px;
                        background: rgba(255, 255, 255, 0.1);
                        transition: all 0.2s ease;
                    `;
                    
                    // Checkbox container for better styling
                    const checkboxContainer = document.createElement('div');
                    checkboxContainer.style.cssText = `
                        display: flex;
                        align-items: center;
                        margin-right: 12px;
                    `;
                    
                    const checkbox = document.createElement('input');
                    checkbox.type = 'checkbox';
                    checkbox.checked = completedFiles.has(file);
                    checkbox.style.cssText = `
                        width: 16px;
                        height: 16px;
                        cursor: pointer;
                        accent-color: #4CAF50;
                    `;
                    
                    checkbox.addEventListener('change', (e) => {
                        if (e.target.checked) {
                            completedFiles.add(file);
                            a.style.textDecoration = 'line-through';
                            a.style.opacity = '0.6';
                        } else {
                            completedFiles.delete(file);
                            a.style.textDecoration = 'none';
                            a.style.opacity = '1';
                        }
                        chrome.storage.local.set({ 'completedFiles': Array.from(completedFiles) });
                    });
                    
                    const a = document.createElement('a');
                    a.href = '#';
                    a.textContent = file;
                    a.style.cssText = `
                        color: inherit;
                        text-decoration: none;
                        flex-grow: 1;
                        font-size: 14px;
                        transition: all 0.2s ease;
                        ${completedFiles.has(file) ? 'text-decoration: line-through; opacity: 0.6;' : ''}
                    `;
                    
                    a.addEventListener('mouseover', () => {
                        if (!completedFiles.has(file)) {
                            li.style.background = 'rgba(255, 255, 255, 0.2)';
                        }
                    });
                    
                    a.addEventListener('mouseout', () => {
                        li.style.background = 'rgba(255, 255, 255, 0.1)';
                    });
                    
                    a.addEventListener('click', (e) => {
                        e.preventDefault();
                        fetchFileContent(file);
                    });
                    
                    checkboxContainer.appendChild(checkbox);
                    li.appendChild(checkboxContainer);
                    li.appendChild(a);
                    filesList.appendChild(li);
                });
    
                // Add reset button if it doesn't exist
                if (!document.getElementById('resetFileListBtn')) {
                    const resetBtn = document.createElement('button');
                    resetBtn.id = 'resetFileListBtn';
                    resetBtn.textContent = 'Reset File List';
                    resetBtn.style.cssText = `
                        margin-top: 12px;
                        padding: 6px 12px;
                        border: none;
                        border-radius: 4px;
                        background: rgba(244, 67, 54, 0.1);
                        color: #F44336;
                        cursor: pointer;
                        font-size: 12px;
                        transition: all 0.2s ease;
                    `;
                    
                    resetBtn.addEventListener('mouseover', () => {
                        resetBtn.style.background = 'rgba(244, 67, 54, 0.2)';
                    });
                    
                    resetBtn.addEventListener('mouseout', () => {
                        resetBtn.style.background = 'rgba(244, 67, 54, 0.1)';
                    });
                    
                    resetBtn.addEventListener('click', function() {
                        if (confirm('Are you sure you want to reset the file list? This will remove all files and their checkmarks.')) {
                            // Clear local storage
                            chrome.storage.local.remove(['persistentFiles', 'completedFiles'], function() {
                                console.log('File list data cleared');
                                
                                // Clear the displayed list
                                const filesList = document.getElementById('addedFilesList');
                                if (filesList) {
                                    filesList.innerHTML = '';
                                }
                                
                                // Reset the server-side list
                                fetch('http://localhost:8080/resetFiles', {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json'
                                    }
                                })
                                .then(response => response.json())
                                .then(data => {
                                    if (data.error) {
                                        console.error('Error resetting files:', data.error);
                                    } else {
                                        console.log('Files reset successfully');
                                        document.getElementById('latestFileName').textContent = 'None';
                                    }
                                })
                                .catch(error => {
                                    console.error('Error resetting files:', error);
                                });
                            });
                        }
                    });
                    
                    document.getElementById('addedFilesSection').appendChild(resetBtn);
                }
            })
            .catch(error => console.error('Error fetching added files:', error));
        });
    }
    
    function fetchFileContent(filename) {
        fetch(`http://localhost:8080/getFileContent?filename=${encodeURIComponent(filename)}`)
        .then(response => response.text())
        .then(content => {
            const contentArea = document.getElementById('fileContent');
            contentArea.textContent = content;
        })
        .catch(error => console.error('Error fetching file content:', error));
    }

    function updateCurrentBranch() {
        fetch('http://localhost:8080/getCurrentBranch')
        .then(response => response.json())
        .then(data => {
            document.getElementById('currentBranch').textContent = data.latestFile || 'Unknown';
        })
        .catch(error => console.error('Error fetching current branch:', error));
    }
    function changeBranch(branchName) {
        fetch(`http://localhost:8080/changeBranch?branch=${encodeURIComponent(branchName)}`, {
            method: 'POST'
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                alert('Error changing branch: ' + data.error);
            } else {
                alert('Successfully changed branch: ' + data.message);
                updateCurrentBranch();
            }
        })
        .catch(error => {
            console.error('Error changing branch:', error);
            alert('Error changing branch: ' + error.message);
        });
    }

    updateCurrentBranch();
    updateAddedFilesList();

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
                if (result.project.social.discord && result.project.social.discord.length > 0) {
                    document.getElementById('discord').value = result.project.social.discord[0].url;
                }
            }
            console.log('Project data loaded:', result.project);
        }
    });

    // Load saved data when popup opens
    loadFormData();

    document.querySelectorAll('input[type="text"]').forEach(input => {
        input.addEventListener('input', saveFormData);
    });

    document.querySelectorAll('input[type="text"]').forEach(input => {
        input.addEventListener('input', saveFormData);
    });

    document.getElementById('changeToSquasherBtn').addEventListener('click', function() {
        changeBranch('squasher');
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

        if (document.getElementById('discord').value) {
            project.social.discord = [{ url: document.getElementById('discord').value }];
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

                // Update the latest file name display
                document.getElementById('latestFileName').textContent = data.latestFile || 'None';

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

    function cleanCommits() {
        fetch('http://localhost:8080/cleanCommits', {
            method: 'POST',
        })
        .then(response => response.json())
        .then(data => {
            console.log('Clean commits result:', data);
            if (data.error) {
                alert('Error cleaning commits: ' + data.error);
            } else {
                alert('Commits cleaned successfully: ' + data.message);
            }
        })
        .catch((error) => {
            console.error('Error cleaning commits:', error);
            alert('Error cleaning commits: ' + error.message);
        });
    }

    // Function to fetch and update the latest file name
    function updateLatestFileName() {
        fetch('http://localhost:8080/getLatestFile')
        .then(response => response.json())
        .then(data => {
            if (data.latestFile) {
                document.getElementById('latestFileName').textContent = data.latestFile;
            }
        })
        .catch(error => console.error('Error fetching latest file name:', error));
    }
    // Update the latest file name when the popup opens
    updateLatestFileName();


    document.getElementById('clearDataBtn').addEventListener('click', function() {
        chrome.storage.sync.remove('project', function() {
            console.log('Project data cleared');
            // Clear all input fields
            document.querySelectorAll('input[type="text"]').forEach(input => {
                input.value = '';
            });
        });
        chrome.storage.local.clear(function() {
            console.log('Local storage cleared');
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
