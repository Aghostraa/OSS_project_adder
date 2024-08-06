console.log('Background script loaded');

chrome.runtime.onInstalled.addListener(function() {
    chrome.contextMenus.create({
        id: "importUrl",
        title: "Import URL to OSS Extension",
        contexts: ["link", "image"]
    });

    chrome.contextMenus.create({
        id: "createDescription",
        title: "Create Project Description",
        contexts: ["selection"]
    });
});

chrome.contextMenus.onClicked.addListener(function(info, tab) {
    if (info.menuItemId === "importUrl") {
        const url = info.linkUrl || info.srcUrl;
        if (url) {
            handleUrl(url);
        }
    } else if (info.menuItemId === "createDescription") {
        if (info.selectionText) {
            createDescription(info.selectionText);
        }
    }
});

function handleUrl(url) {
    copyToClipboard(url);
    const category = categorizeUrl(url);
    saveUrl(url, category);
}

function copyToClipboard(text) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    console.log('Text copied to clipboard:', text);
    showNotification('Text Copied', 'The selected text has been copied to the clipboard.');
}

function categorizeUrl(url) {
    if (url.includes('github.com')) return 'github';
    if (url.includes('twitter.com') || url.includes('t.co')) return 'twitter';
    if (url.includes('t.me')) return 'telegram';
    if (url.includes('mirror.xyz')) return 'mirror';
    return 'website';
}

function saveUrl(url, category) {
    chrome.storage.local.set({ [category]: url, selectedUrl: url }, function() {
        console.log('URL saved:', url, 'Category:', category);
    });
}

function showNotification(title, message) {
    chrome.notifications.create({
        type: 'basic',
        iconUrl: 'icons/icon48.png',
        title: title,
        message: message
    }, function(notificationId) {
        console.log('Notification created with ID:', notificationId);
    });
}
function createDescription(selectedText) {
    const prompt = "Create a concise, neutral, no marketing, 1 liner description of the project from the following intro:\n\n" + selectedText;
    copyToClipboard(prompt);
    showNotification('Description Prompt Copied', 'The description prompt has been copied to your clipboard.');
}