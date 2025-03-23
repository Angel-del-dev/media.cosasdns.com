const do_login = async e => {
    e.preventDefault();

    const User = document.querySelector('[name="username"]').value.trim();
    const Password = document.querySelector('[name="password"]').value.trim();
    const ErrorNode = document.getElementById('ErrorNode');
    ErrorNode.classList.add('d-none');
    if(User === '' || Password === '') {
        ErrorNode.classList.remove('d-none');
        ErrorNode.textContent = 'Invalid credentials';
        return ;
    }
};

document.getElementById('SignIn').addEventListener('click', do_login);