fetch('/api/items')
  .then(r => r.json())
  .then(items => {
    const list = document.getElementById('list');
    items.forEach(item => {
      const li = document.createElement('li');
      li.textContent = item.name;
      list.appendChild(li);
    });
  })
  .catch(err => console.error('Error:', err));