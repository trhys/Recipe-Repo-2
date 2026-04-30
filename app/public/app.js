// Page load
document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  const refresher = localStorage.getItem('refresh_token');
});

// Print shopping list
async function printList(shoppingListId) {
	const token = localStorage.getItem('token');

	try {
		const url = `/shoppinglists/${shoppingListId}/print`
		const body = {
			method: 'GET',
			headers: {
				'Authorization': `Bearer ${token}`,
				'Accept': 'text/html'
			},
		}
		let response = await fetch(url, body);

		if (response.ok) {
			const newHTML = await response.text();

			document.open();
			document.write(newHTML);
			document.close();
			window.history.pushState({}, '', `/shoppinglists/${shoppingListId}/print`);
			return
		}

		const refreshRes = await fetch('/api/tokens/refresh', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${localStorage.getItem('refresh_token')}`,
			},
		});

		if (!refreshRes.ok) throw new Error("Session expired");
		const data = await refreshRes.json();
		if (data.token) {
			localStorage.setItem('token', data.token);
		}

		response = await fetch(url, body);

		if (!response.ok) {
			const data = await response.json();
			throw new Error(`Failed to renew token: ${data.error}`);
		}
	} catch (error) {
		console.error('Error:', error);
	}
}

printBtn = document.getElementById('print-shopping-list')
if (printBtn) {
	printBtn.addEventListener('click', async (event) => {
		event.preventDefault();
		await printList(printBtn.value);
	});
}
