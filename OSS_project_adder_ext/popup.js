document.addEventListener('DOMContentLoaded', function() {
    // Retrieve stored form data
    chrome.storage.sync.get(['project'], function(result) {
        if (result.project) {
            document.getElementById('name').value = result.project.name;
            document.getElementById('displayName').value = result.project.displayName;
            document.getElementById('description').value = result.project.description;
            if (result.project.websites.length > 0) {
                document.getElementById('website').value = result.project.websites[0];
            }
            if (result.project.github.length > 0) {
                document.getElementById('github').value = result.project.github[0];
            }
            if (result.project.social) {
                if (result.project.social.twitter.length > 0) {
                    document.getElementById('twitter').value = result.project.social.twitter[0];
                }
                if (result.project.social.telegram.length > 0) {
                    document.getElementById('telegram').value = result.project.social.telegram[0];
                }
                if (result.project.social.mirror.length > 0) {
                    document.getElementById('mirror').value = result.project.social.mirror[0];
                }
            }
        }
    });

    document.getElementById('projectForm').addEventListener('submit', function(event) {
        event.preventDefault();

        const project = {
            name: document.getElementById('name').value,
            displayName: document.getElementById('displayName').value,
            description: document.getElementById('description').value,
            websites: document.getElementById('website').value ? [document.getElementById('website').value] : [],
            github: document.getElementById('github').value ? [document.getElementById('github').value] : [],
            social: {
                twitter: document.getElementById('twitter').value ? [document.getElementById('twitter').value] : [],
                telegram: document.getElementById('telegram').value ? [document.getElementById('telegram').value] : [],
                mirror: document.getElementById('mirror').value ? [document.getElementById('mirror').value] : []
            }
        };

        // Save form data to storage
        chrome.storage.sync.set({ project: project }, function() {
            console.log('Project data saved');
        });

        fetch('http://localhost:8080/createProject', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(project)
        })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            alert('Project created successfully!');
        })
        .catch((error) => {
            console.error('Error:', error);
            alert('Error creating project');
        });
    });
});