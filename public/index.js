// Elements
const uploadInput = document.querySelector('#hidden-upload');
const adminPasswordInput = document.querySelector('#admin-password');
const passwordInput = document.querySelector('#password')
const errorDiv = document.querySelector('#error-div');
const errorTitle = document.querySelector('#error-title');
const errorMessage = document.querySelector('#error-message');
const displayFilesButton = document.querySelector('#display-files');

// Event Handlers
uploadInput.addEventListener('change', upload)
displayFilesButton.addEventListener('click', getFiles)

// Functions
function displayError(title, message) {
	if (!title) errorDiv.style.display = 'none';
	else {
		errorDiv.style.display = 'block';
		errorTitle.innerText = title;
		if (message) errorMessage.innerText = message;
	}
}

function getFiles() {
	let adminPassword = adminPasswordInput.value || 'invalid';
	document.querySelector('#files-response-container').innerHTML = '';

	// Get files from api
	fetch('/api/files/' + adminPassword).then(response => {
		displayError()

		let status = response.status;

		if (status === 200) {
			response.json().then(data => {
				if (data instanceof Array) {
					for (let i = 0; i < data.length; i++) {
						console.log(data[i])
						appendFile("files-response-container", {
							filename: data[i].Name,
							url: `${window.location.origin}/api/${data[i].Name}`,
							deletionUrl: `${window.location.origin}/api/${data[i].Name}/delete/${data[i].EditKey}`,
							subtext: `${data[i].Views} views`
						})
					}
				}
			});
		}
		else if (status === 401) {
			displayError('Unauthorized', 'The admin password you entered was invalid.');
		}

	});
}

function appendFile(containerId, file) {
	let container = document.getElementById(containerId);
	let fileDiv = document.createElement('div');
	fileDiv.className = 'item';
	fileDiv.innerHTML = `<div class="right floated content"><div class="ui button inverted" onclick="navigator.clipboard.writeText('${file.url}')">Copy URL</div><a href="${file.deletionUrl}" target="_blank"><div class="ui button red inverted">Delete</div></a></div><i class="file outline icon"></i><div class="content">${file.filename}<br><span class="subtext">${file.subtext}</span></div>`
	container.appendChild(fileDiv);
}

function upload() {
	let files = uploadInput.files;
	let password = passwordInput.value;

	for (let i = 0; i < files.length; i++) {
		let file = files[i];

		// Upload file to api
		let formData = new FormData();
		formData.append('file', file);
		formData.append('password', password);
		fetch('/api/upload', {
			method: 'POST',
			body: formData
		}).then(response => {
			displayError()

			let status = response.status;
			if (status === 200) {
				response.json().then(data => {
					appendFile("upload-response-container", { subtext: `${data.size} bytes`, ...data })
				});
			} else if (status === 401) {
				displayError('Unauthorized', 'The password you entered was invalid.');
			}
		});

	}

	uploadInput.value = '';
}