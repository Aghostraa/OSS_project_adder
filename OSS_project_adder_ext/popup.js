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
