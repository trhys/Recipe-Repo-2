export async function fetchDbOptions() {
  try {
    const response = await fetch('/api/ingredients');
    const items = await response.json();
    populateList(items.ingredients);
      } catch (err) {
    console.error("Failed to load database options:", err);
  }
}

async function getUnits(ingredientId, units) {
  try {
    const url = '/api/units'
    const body = {
            method: 'POST',
            body: JSON.stringify({ ingredient_id: ingredientId }),
    }
    const unitsResponse = await fetch(url, body);
    const data = await unitsResponse.json();
    populateUnits(data.units, units);
      } catch (err) {
    console.error("Failed to load database options:", err);
  }
}

function populateList(items) {
  const list = document.getElementById('ingredient-list');
  list.innerHTML = items.map(item => `
          <option value="${item.name}" data-id="${item.id}"></option>
  `).join('');
}

function populateUnits(units, row) {
        row.innerHTML = units.map(item => `
                <option value="${item.name}" data-id="${crypto.randomUUID()}">${item.name}</option>
        `).join('');
}

function filterOptions() {
  const input = document.getElementById('dbSearch').value.toLowerCase();
  const listItems = document.querySelectorAll('#ingredient-list option');

  listItems.forEach(option => {
    const text = option.textContent.toLowerCase();
    option.style.display = text.includes(input) ? "block" : "none";
  });
}

const wrapper = document.getElementById('ingredients-wrapper');
wrapper.addEventListener('input', function(e) {
  if (e.target.classList.contains('db-search-input')) {
    const inputValue = e.target.value;
    const row = e.target.closest('.ingredient-row');
    
    const unitSelect = row.querySelector('.ingredient-units'); 
    const list = document.getElementById('ingredient-list');
    
    const selectedOption = Array.from(list.options).find(option => option.value === inputValue);

    if (selectedOption) {
      const id = selectedOption.getAttribute('data-id');
      
      const hiddenInput = row.querySelector('.ingredient_uuid');
      if (hiddenInput) hiddenInput.value = id;

      getUnits(id, unitSelect);
    }
  }
});


// Form listeners
const addIngBtn = document.getElementById('add-ingredient');
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

const removeIngBtn = document.getElementById('ingredients-wrapper');
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

const inWrapper = document.getElementById('ingredients-wrapper')
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
