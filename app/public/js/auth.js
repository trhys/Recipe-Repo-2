import {
    login,
    signup
} from './api.js'

export const auth = {
    async getSession() {
        const data = await login();
        this.saveSession(data);
    },

    saveSession(data) {
        localStorage.setItem('token', data.token);
        localStorage.setItem('refresh_token', data.refresh_token);
        localStorage.setItem('user_id', data.id);
    },

    logout() {
        localStorage.clear();
        window.location.href = '/login';
    },

    isLoggedIn() {
        return !!localStorage.getItem('token');
    },

    getHeaders() {
        const token = localStorage.getItem('token');
        return {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        };
    },

    updateNav() {
	const loggedIn = this.isLoggedIn()

	document.querySelectorAll('[data-show-if]').forEach(element => {
		const show = element.getAttribute('data-show-if');

		if (show === 'logged-in') {
			element.style.display = loggedIn ? 'block' : 'none';
		} else if (show === 'logged-out') {
			element.style.display = loggedIn ? 'none' : 'block';
		}
	});
    },
};
