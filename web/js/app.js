let personas = [];
let currentScore = null;

async function init() {
    await loadPersonas();
}

async function loadPersonas() {
    try {
        const response = await fetch('/api/personas');
        personas = await response.json();
        renderPersonas();
    } catch (error) {
        console.error('Failed to load personas:', error);
    }
}

function renderPersonas() {
    const grid = document.getElementById('personasGrid');
    const avatars = ['👩‍🏫', '👨‍💻', '🧑‍🔧', '👩‍🎨', '👵'];
    
    grid.innerHTML = personas.map((p, i) => `
        <div class="persona-card" onclick="getScore('${p.id}')">
            <div class="persona-avatar">${avatars[i]}</div>
            <div class="persona-name">${p.name}</div>
            <div class="persona-occupation">${p.occupation}</div>
        </div>
    `).join('');
}

async function getScore(personaId) {
    try {
        const response = await fetch(`/api/personas/${personaId}/score`);
        currentScore = await response.json();
        showScoreReveal();
    } catch (error) {
        console.error('Failed to get score:', error);
    }
}

function showLanding() {
    document.querySelectorAll('.view').forEach(v => v.classList.add('hidden'));
    document.getElementById('landing').classList.remove('hidden');
    document.getElementById('landing').classList.add('active');
}

function showCustomInput() {
    document.querySelectorAll('.view').forEach(v => v.classList.add('hidden'));
    document.getElementById('customInput').classList.remove('hidden');
    document.getElementById('customInput').classList.add('active');
}

function showScoreReveal() {
    document.querySelectorAll('.view').forEach(v => v.classList.add('hidden'));
    document.getElementById('scoreReveal').classList.remove('hidden');
    document.getElementById('scoreReveal').classList.add('active');
    
    animateScore(currentScore);
}

function animateScore(score) {
    document.getElementById('scoreNumber').textContent = '0';
    document.getElementById('scoreBand').textContent = score.scoreBand;
    document.getElementById('patternValue').textContent = score.pattern_detected.replace(/_/g, ' ');
    
    const circumference = 2 * Math.PI * 90;
    const progress = document.querySelector('.dial-progress');
    const offset = circumference - (score.total_score / 900) * circumference;
    
    setTimeout(() => {
        progress.style.strokeDashoffset = offset;
    }, 100);
    
    let current = 0;
    const target = score.total_score;
    const duration = 1500;
    const increment = target / (duration / 16);
    
    const counter = setInterval(() => {
        current += increment;
        if (current >= target) {
            current = target;
            clearInterval(counter);
        }
        document.getElementById('scoreNumber').textContent = Math.round(current);
    }, 16);
    
    renderRadarChart(score.components);
    renderInsights(score.insights);
    renderProducts(score.credit_products);
    renderImprovements(score.improvements);
}

function renderInsights(insights) {
    const list = document.getElementById('insightsList');
    list.innerHTML = insights.map(i => `<p>${i}</p>`).join('');
}

function renderProducts(products) {
    const grid = document.getElementById('productsGrid');
    grid.innerHTML = products.map(p => `
        <div class="product-card">
            <h4>${p.name}</h4>
            <div class="product-type">${p.type}</div>
            <div class="product-details">
                Limit: ₹${p.limit.toLocaleString()}<br>
                Interest: ${p.interest}%
            </div>
        </div>
    `).join('');
}

function renderImprovements(improvements) {
    const list = document.getElementById('improvementsList');
    list.innerHTML = improvements.map(i => `
        <div class="improvement-card">
            <span class="improvement-action">${i.action}</span>
            <span class="improvement-points">+${i.points_delta} pts</span>
        </div>
    `).join('');
}

document.getElementById('customForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const name = document.getElementById('inputName').value;
    const age = parseInt(document.getElementById('inputAge').value);
    const city = document.getElementById('inputCity').value;
    const depositsJson = document.getElementById('inputDeposits').value;
    
    let deposits = [];
    try {
        deposits = JSON.parse(depositsJson);
    } catch (err) {
        alert('Invalid JSON format for deposits');
        return;
    }
    
    const history = {
        user_id: 'custom',
        name: name,
        age: age,
        city: city,
        deposits: deposits
    };
    
    try {
        const response = await fetch('/api/score', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(history)
        });
        currentScore = await response.json();
        showScoreReveal();
    } catch (error) {
        console.error('Failed to calculate score:', error);
    }
});

init();