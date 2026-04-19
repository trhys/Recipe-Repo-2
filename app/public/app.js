// Page load

document.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  const refresher = localStorage.getItem('refresh_token');

  loginBtn = document.getElementById('login-btn')
  signupBtn = document.getElementById('signup-btn')
  logoutBtn = document.getElementById('logout-button')
  recipeCreator = document.getElementById('recipe-creator')
  createNew = document.getElementById('create-new')
  recipeList = document.getElementById('recipes-section')

  if (loginBtn) {
	  if (token) {
		loginBtn.style.display = 'none';
	  } else {
		loginBtn.style.display = 'block';
	  }
  }
  if (signupBtn) {
	if (token) {
		signupBtn.style.display = 'none';
	} else {
		signupBtn.style.display = 'block';
	}
  }
  if (logoutBtn) {
	if (token) {
		logoutBtn.style.display = 'block';
	} else {
		logoutBtn.style.display = 'none';
	}
  }
  if (recipeCreator) {
	if (token) {
		recipeCreator.style.display = 'block';
	} else {
		recipeCreator.style.display = 'none';
	}
  }
  if (createNew) {
	if (token) {
		createNew.style.display = 'block';
	} else {
		createNew.style.display = 'none';
	}
  }
  if (recipeList) {
	  await getRecipes();
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
      localStorage.setItem('refresh_token', data.refresh_token);
      localStorage.setItem('user_id', data.id);
      window.location.href='/';
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
    window.location.href='/login';  
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}

function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('user_id');
  window.location.href='/';
}

// Recipe view

async function getRecipes() {
	try {
    const res = await fetch('/api/recipes', {
      method: 'GET',
    });
    if (!res.ok) {
      const data = await res.json();
      throw new Error(`Failed to get recipe list. Error: ${data.error}`);
    }

    const data = await res.json();
    const recipes = data.recipes;
    const recipeList = document.getElementById('recipes');
    recipeList.innerHTML = '';
    if (recipes == null) return;
    for (const recipe of recipes) {
      const listItem = document.createElement('li');
      const displayDate = new Date(recipe.created_at).toLocaleDateString();
      listItem.innerHTML = `
      <div class="card-content">
	<a href="/recipes/${recipe.id}">
      	<h4>${recipe.title}</h4></a>
	<p>Created: ${displayDate}</p>
	<a href="/users/${recipe.user_id}">
	<h5>Author: ${recipe.author}</h5></a>
	<br /><div class="content-card">
		<img src=${recipe.image_url} class="card-image">
	</div>
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

inWrapper = document.getElementById('ingredients-wrapper')

if (inWrapper) {
	inWrapper.addEventListener('input', function(e) {
	  if (e.target.classList.contains('db-search-input')) {
	    const row = e.target.closest('.ingredient-row');
	    const hiddenInput = row.querySelector('.ingredient_uuid');
	    const datalist = document.getElementById('ingredient-list');
	    
	    const selectedOption = Array.from(datalist.options).find(opt => opt.value === e.target.value);

	    if (selectedOption) {
	      hiddenInput.value = selectedOption.getAttribute('data-id');
	    } else {
	      hiddenInput.value = ""; 
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
				ingredients: [],
				description: document.getElementById('author-description').value
			};

			const rows = document.querySelectorAll('.ingredient-row')

			rows.forEach(row => {
				const ingredient = {
					id: row.querySelector('input[name="ingredient_id"]').value,
					quantity: parseFloat(row.querySelector('input[name="quantity"]').value),
					unit: row.querySelector('select[name="unit"]').value
				};
				recipeData.ingredients.push(ingredient);
			});

			alert(`Sending: ${JSON.stringify(recipeData)}`);

			const formData = new FormData();
			formData.append("payload", JSON.stringify(recipeData))
			formData.append("image", recipeCreator.recipe_pic.files[0])

			const url = '/api/new_recipe'
			const reqBody = {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('token')}`,
				},
				body: formData,
			}

			let res = await fetch(url, reqBody);
			if (res.status !== 401) {
				window.location.href = '/';
				return res;
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

			res = await fetch(url, reqBody);

			if (!res.ok) {
				const data = await res.json();
				throw new Error(`Failed to create recipe: ${data.error}`);
			}
			window.location.href='/';
		} catch (error) {
			alert(`Error: ${error.message}`);
		}
	});
}

addIngredient = document.getElementById('add-ingredient-panel');
if (addIngredient) {
	addIngredient.addEventListener('submit', async (event) => {
		event.preventDefault();
		try {
			const ingredientData = {
				name: document.getElementById('ingredient-name').value,
			}

			const formData = new FormData();
			formData.append("payload", JSON.stringify(ingredientData));

			const url = '/api/admin/new_ingredient'
			const reqBody = {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('token')}`,
				},
				body: formData,
			}

			let res = await fetch(url, reqBody);
			if (res.status !== 401) {
				alert('Success')
				window.location.href = '/';
				return res;
			}

			const refreshRes = await fetch('/api/refresh', {
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

			res = await fetch(url, reqBody);

			if (!res.ok) {
				const data = await res.json();
				throw new Error(`Failed to create ingredient: ${data.error}`);
			}
			window.location.href='/';
		} catch (error) {
			alert(`Error: ${error.message}`);
		}
	});
}
