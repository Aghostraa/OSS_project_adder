console.log('Content script loaded');

document.addEventListener('mouseup', function(event) {
    console.log('Mouse up event detected');
    let selectedText = window.getSelection().toString().trim();
    console.log('Selected text:', selectedText);
    if (selectedText) {
        chrome.runtime.sendMessage({
            action: "textSelected",
            text: selectedText
        }, function(response) {
            console.log('Response from background script:', response);
        });
    }
});