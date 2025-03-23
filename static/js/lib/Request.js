export const Request = async ({
    route = null,
    method = 'POST',
    data = {}
}) => {
    if(route === null) throw new Exception("Route not provided");

    const params = {
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/x-www-form-urlencoded"
        },
        method
    }
    if(!['GET'].includes(method.toUpperCase())) params.body = new URLSearchParams(data);
    if(localStorage.getItem('auth') !== null) params.headers.Authorization = `Bearer ${localStorage.getItem('auth')}`

    return await fetch(route, params)
        .then(r =>  r.ok || r.status === 400 ? r.json() : { message: r.statusText })
        .then(r => {
            if(r.data !== undefined) {
                const data = JSON.parse(r.data);
                if(data.token !== undefined) {
                    localStorage.setItem('auth', data.token);
                    delete data.token;
                }
                data.message = '';
                return data;
            }
            return r;
        });
};