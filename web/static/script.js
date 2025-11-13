document.getElementById('shorten-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const longUrl = document.getElementById('long-url').value;
    const resultDiv = document.getElementById('result');
    const errorDiv = document.getElementById('error');
    const shortUrlLink = document.getElementById('short-url');
    
    // Hide previous results
    resultDiv.classList.add('hidden');
    errorDiv.classList.add('hidden');
    
    try {
        const response = await fetch('/api/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ long_url: longUrl }),
        });
        
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }
        
        const data = await response.json();
        
        // Display result
        shortUrlLink.textContent = data.short_url;
        shortUrlLink.href = data.short_url;
        resultDiv.classList.remove('hidden');
        
        // Clear input
        document.getElementById('long-url').value = '';
        
    } catch (error) {
        errorDiv.textContent = 'Error: ' + error.message;
        errorDiv.classList.remove('hidden');
    }
});