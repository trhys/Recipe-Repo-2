import { auth } from './auth.js'

export async function loadNav() {
	const response = await fetch('/components/navbar.html');
	const html = await response.text();

	const navbar = document.getElementById('navbar');
	if (navbar) {
		navbar.innerHTML = html;

		auth.updateNav();

		const shopListBtn = document.getElementById('shopping-lists-btn')
		if (shopListBtn) {
		    shopListBtn.addEventListener('click', async (event) => {
			event.preventDefault();
			await viewShoppingLists();
		    });
		}

		const logoutBtn = document.getElementById('logout-button')
		if (logoutBtn) {
			logoutBtn.addEventListener('click', async (event) => {
				event.preventDefault();
				await auth.logout();
			});
		}
	}
}
