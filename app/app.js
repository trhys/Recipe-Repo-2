// Page load

document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');

  loginBtn = document.getElementById('login-btn')
  signupBtn = document.getElementById('signup-btn')
  logoutBtn = document.getElementById('logout-button')
  recipeCreator = document.getElementById('recipe-creator')

  if (token) {
	  if (loginBtn) {
    		loginBtn.style.display = 'none';
	  }
	  if (signupBtn) {
    		signupBtn.style.display = 'none';
	  }
	  if (logoutBtn) {
    		logoutBtn.style.display = 'block';
	  }
	  if (recipeCreator) {
    		recipeCreator.style.display = 'block';
		await getRecipes();
	  }
  } else {
	  if (loginBtn) {
    		loginBtn.style.display = 'block';
	  }
	  if (signupBtn) {
    		signupBtn.style.display = 'block';
	  }
	  if (logoutBtn) {
    		logoutBtn.style.display = 'none';
	  }
	  if (recipeCreator) {
    		recipeCreator.style.display = 'none';
		await getRecipes();
	  }
  }
});

loginForm = document.getElementById('login-form');
if (loginForm) {
	loginForm.addEventListener('submit', async (event) => {
  		event.preventDefault();
  		await login();
	});
}

signupForm = document.getElementById('signup-form');
if (signupForm) {
	signupForm.addEventListener('submit', async (event) => {
		event.preventDefault();
		await signup();
	});
}

// Login/Signup

async function login() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to login: ${data.error}`);
    }

    if (data.token) {
      localStorage.setItem('token', data.token);
      localStorage.setItem('user_id', data.id);
      window.location.href='/app';
    } else {
      alert('Login failed. Please check your credentials.');
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

async function signup() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const name = document.getElementById('username').value;

  try {
    const res = await fetch('/api/new_user', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, name }),
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to create user: ${data.error}`);
    }
    window.location.href='/app/login';  
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  window.location.href='/app';
}

// Recipe view

async function getRecipes() {
	try {
    const res = await fetch('/api/recipes', {
      method: 'GET',
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to get videos. Error: ${data.error}`);
    }

    const data = await res.json();
    const recipes = data.recipes;
    const recipeList = document.getElementById('recipes');
    recipeList.innerHTML = '';
    for (const recipe of recipes) {
      const listItem = document.createElement('li');
      const displayDate = new Date(recipe.created_at).toLocaleDateString();
      listItem.innerHTML = `
      <div class="card-content">
      	<h4>${recipe.title}</h4>
	<p>Created: ${displayDate}</p>
	<h5>Author: ${recipe.author}</h5>
      </div>`;
      recipeList.appendChild(listItem);
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

// Add Recipe

addIngBtn = document.getElementById('add-ingredient');
if (addIngBtn) {
	addIngBtn.addEventListener('click', function(event) {
		event.preventDefault();

		var wrapper = document.getElementById('ingredients-wrapper');
		
		const firstRow = document.querySelector('.ingredient-row');
		const newRow = firstRow.cloneNode(true);
		const inputs = newRow.querySelectorAll('input');
		inputs.forEach(input => input.value = '');
		wrapper.appendChild(newRow);
	});
}

removeIngBtn = document.getElementById('ingredients-wrapper');
if (removeIngBtn) {
	removeIngBtn.addEventListener('click', function(event) {
		if (event.target.classList.contains('remove-btn')) {
			event.preventDefault();

			const rows = document.querySelectorAll('.ingredient-row');
			if (rows.length > 1) {
				event.target.closest('.ingredient-row').remove();
			} else {
				alert("Recipe must contain at least one ingredient!");
			}
		}
	});
}

recipeCreator = document.getElementById('recipe-creator');
if (recipeCreator) {
	recipeCreator.addEventListener('submit', async (event) => {
		event.preventDefault();
		try {
			const recipeData = {
				title: document.getElementById('recipe-title').value,
				user_id: localStorage.getItem('user_id'),
				ingredients: []
			};

			const rows = document.querySelectorAll('.ingredient-row')

			rows.forEach(row => {
				const ingredient = {
					name: row.querySelector('input[name="name"]').value,
					quantity: parseFloat(row.querySelector('input[name="quantity"]').value),
					unit: row.querySelector('input[name="unit"]').value
				};
				recipeData.ingredients.push(ingredient);
			});

			const res = await fetch('/api/new_recipe', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${localStorage.getItem('token')}`,
				},
				body: JSON.stringify(recipeData),
			});
			console.log("Sending request:", JSON.stringify(recipeData, null, 2));
			if (!res.ok) {
				const data = await res.json();
				throw new Error(`Failed to create recipe: ${data.error}`);
			}
			window.location.href='/app';
		} catch (error) {
			alert(`Error: ${error.message}`);
		}
	});
}
