document.addEventListener('DOMContentLoaded', function() {
    console.log('Popup loaded');

    function saveFormData() {
        const project = {
            name: document.getElementById('name').value,
            displayName: document.getElementById('displayName').value,
            description: document.getElementById('description').value,
            websites: document.getElementById('website').value ? [{ url: document.getElementById('website').value }] : [],
            github: document.getElementById('github').value ? [{ url: document.getElementById('github').value }] : [],
            social: {
                twitter: document.getElementById('twitter').value ? [{ url: document.getElementById('twitter').value }] : [],
                telegram: document.getElementById('telegram').value ? [{ url: document.getElementById('telegram').value }] : [],
                mirror: document.getElementById('mirror').value ? [{ url: document.getElementById('mirror').value }] : []
            }
        };
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
                if (result.project.social.twitter.length > 0) {
                    document.getElementById('twitter').value = result.project.social.twitter[0].url;
                }
                if (result.project.social.telegram.length > 0) {
                    document.getElementById('telegram').value = result.project.social.telegram[0].url;
                }
                if (result.project.social.mirror.length > 0) {
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
            social: {
                twitter: document.getElementById('twitter').value ? [{ url: document.getElementById('twitter').value }] : [],
                telegram: document.getElementById('telegram').value ? [{ url: document.getElementById('telegram').value }] : [],
                mirror: document.getElementById('mirror').value ? [{ url: document.getElementById('mirror').value }] : []
            }
        };

        chrome.storage.sync.set({ project: project }, function() {
            console.log('Project data saved before submission:', project);
        });

        console.log('Sending request to server with project data:', project);

        fetch('http://localhost:8080/createProject', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(project)
        })
        .then(response => {
            console.log('Received response:', response);
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(data => {
            console.log('Success:', data);
            alert('Project created successfully!');
            chrome.storage.sync.remove('project', function() {
                console.log('Project data cleared from storage');
            });
            document.getElementById('projectForm').reset();
        })
        .catch((error) => {
            console.error('Error:', error);
            alert('Error creating project: ' + error.message);
        });
    });
});
