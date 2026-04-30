import { auth } from './auth.js'
import { 
	getShoppingListsHTML,
	getShoppingListsJSON,
	postAddtoList 
} from './api.js'

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
			await getShoppingListsHTML();
			return
		    });
		}

		const logoutBtn = document.getElementById('logout-button')
		if (logoutBtn) {
			logoutBtn.addEventListener('click', async (event) => {
				event.preventDefault();
				await auth.logout();
				return
			});
		}
	}
}

export async function addToListModal(recipeID) {
	const response = await fetch('/components/addPopup.html');
	const html = await response.text();

	const modal = document.getElementById('modal');
	if (modal) {
		modal.innerHTML = html;

		modal.showModal();
		const select = modal.querySelector('#shopping-lists');

		try {
			const data = await getShoppingListsJSON();

			data.shopping_lists.forEach(list => {
				const opt = document.createElement('option');
				opt.value = list.id;
				opt.textContent = list.name;
				select.appendChild(opt);
			});
		} catch (error) {
			alert(`Error: ${error.message}`);
		}

		const closeBtn = modal.querySelector('#closePopup')
		if (closeBtn) {
			closeBtn.onclick = () => {
			    modal.close();
			};
		}

		const form = modal.querySelector('#addtolistform');
		if (form) {
			form.addEventListener('submit', async (event) => {
				event.preventDefault();
				const shoppingListId = select.value;
				await postAddtoList(shoppingListId, recipeID);
				modal.close();
			});
		}
	}
}
