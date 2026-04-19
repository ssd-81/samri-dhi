let personas = [];
let currentScore = null;
let currentPersonaId = null;
let customDeposits = [];
let wizardStep = 1;
let radarChart = null;
let timelineChart = null;
let lastScrollPosition = 0;

const avatars = [
    '<svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="38" fill="#e8c4a0"/><circle cx="40" cy="32" r="18" fill="#2d1f14"/><path d="M22 60 Q40 75 58 60" fill="#2d1f14"/><circle cx="32" cy="38" r="2" fill="#2d1f14"/><circle cx="48" cy="38" r="2" fill="#2d1f14"/><path d="M35 45 Q40 48 45 45" fill="none" stroke="#2d1f14" stroke-width="1.5"/></svg>',
    '<svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="38" fill="#d4a574"/><circle cx="40" cy="32" r="18" fill="#1a1a1a"/><path d="M22 60 Q40 75 58 60" fill="#1a1a1a"/><rect x="30" y="30" width="20" height="8" fill="#2563eb" rx="4"/></svg>',
    '<svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="38" fill="#c9a067"/><circle cx="40" cy="32" r="18" fill="#3d2b1f"/><path d="M22 60 Q40 75 58 60" fill="#3d2b1f"/><circle cx="32" cy="38" r="2" fill="#3d2b1f"/><circle cx="48" cy="38" r="2" fill="#3d2b1f"/><path d="M35 44 Q40 47 45 44" fill="none" stroke="#3d2b1f" stroke-width="1.5"/></svg>',
    '<svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="38" fill="#e8c4a0"/><circle cx="40" cy="30" r="18" fill="#4a3728"/><path d="M22 58 Q40 72 58 58" fill="#4a3728"/><circle cx="32" cy="36" r="2" fill="#4a3728"/><circle cx="48" cy="36" r="2" fill="#4a3728"/><path d="M35 42 Q40 45 45 42" fill="none" stroke="#4a3728" stroke-width="1.5"/></svg>',
    '<svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="38" fill="#d4a574"/><circle cx="40" cy="34" r="18" fill="#8b7355"/><path d="M22 62 Q40 75 58 62" fill="#8b7355"/><circle cx="32" cy="40" r="2" fill="#8b7355"/><circle cx="48" cy="40" r="2" fill="#8b7355"/><path d="M35 46 Q40 49 45 46" fill="none" stroke="#8b7355" stroke-width="1.5"/></svg>'
];

const personaTaglines = {
    priya: "9 years of perfect saving. Zero credit history.",
    vikram: "8 banks, rate-optimized, aggressive saver.",
    ramesh: "COVID broke his streak. He's coming back stronger.",
    anita: "She saves consistently... then breaks her FDs for cash flow.",
    sunita: "15 years at one bank. ₹13L in FDs. Still invisible to lenders."
};

const personaStories = {
    priya: {
        narrative: "Priya has created an FD every single year since 2016. She held every one to maturity. She's completed two Recurring Deposits without missing an installment.",
        hook: "Her CIBIL score? Not Available.",
        punchline: "Traditional credit scoring can't see her. Samridhi can."
    },
    vikram: {
        narrative: "Vikram optimizes everything. 8 different banks, staggered maturities, rate-shopping across institutions. He has more financial planning in his pinky than most loan applicants.",
        hook: "He's done more financial planning than most loan applicants.",
        punchline: "His FD behavior proves he'd be a model borrower."
    },
    ramesh: {
        narrative: "Ramesh had a good streak going — annual FDs from 2017 to 2019. Then COVID hit his shop and he had to make hard choices.",
        hook: "He broke two FDs to survive. The system would penalize him for it.",
        punchline: "Samridhi sees more — his last 3 deposits show a powerful comeback."
    },
    anita: {
        narrative: "Anita is a consistent saver — she creates FDs regularly. But she breaks them just as often for emergency cash needs.",
        hook: "She's not irresponsible. She lacks liquidity. She uses her FDs as an emergency fund.",
        punchline: "Give her a small credit line, and she'll stop breaking FDs entirely."
    },
    sunita: {
        narrative: "Sunita retired 5 years ago. She's been saving at SBI — and only SBI — for 15 years with flawless discipline.",
        hook: "₹13 lakh in deposits. All at one bank. All in the same tenure.",
        punchline: "Her discipline is flawless. Her diversification could improve — and we'll show her how."
    }
};

const patternNarratives = {
    DISCIPLINED_OPTIMIZER: "Outstanding! You represent the ideal FD customer — consistent, disciplined, and strategically smart. You're primed for premium credit products with favorable terms.",
    LOYAL_SINGLE_BANK: "Your savings discipline is excellent — years of perfect maturity completion is rare. But keeping everything at one bank means you might be missing better rates elsewhere and exceeding DICGC insurance limits.",
    LIQUIDITY_GAP_SAVER: "You're a consistent saver who occasionally breaks FDs for cash flow. You're not irresponsible — you just need better liquidity options.",
    RECOVERING_SAVER: "You've shown resilience by recovering from past financial challenges. Your growth trajectory is positive — keep building on this foundation.",
    STANDARD_SAVER: "Your FD habits show a solid foundation. Focus on building consistency and discipline to unlock better credit opportunities."
};

async function init() {
    await loadPersonas();
    setupScrollObserver();
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
    
    grid.innerHTML = personas.map((p, i) => `
        <div class="persona-card" onclick="showProfile('${p.id}')">
            <div class="persona-avatar">${avatars[i]}</div>
            <div class="persona-name">${p.name}</div>
            <div class="persona-occupation">${p.occupation}</div>
            <div class="persona-location">${p.age}, ${p.city}</div>
            <div class="persona-tagline">"${personaTaglines[p.id]}"</div>
            <div class="persona-stat">${p.deposits.length} FDs • ${getYearsActive(p.deposits)} years</div>
        </div>
    `).join('');
}

function getYearsActive(deposits) {
    const years = new Set();
    deposits.forEach(d => {
        if (d.start_date && d.start_date.length >= 4) {
            years.add(d.start_date.substring(0, 4));
        }
    });
    return Math.max(years.size, 1);
}

function showProfile(personaId) {
    lastScrollPosition = window.scrollY;
    currentPersonaId = personaId;
    const persona = personas.find(p => p.id === personaId);
    if (!persona) return;

    const avatarMap = { priya: 0, vikram: 1, ramesh: 2, anita: 3, sunita: 4 };
    const story = personaStories[personaId];

    document.getElementById('profileAvatar').innerHTML = avatars[avatarMap[personaId]];
    document.getElementById('profileName').textContent = persona.name;
    document.getElementById('profileMeta').textContent = `${persona.occupation} • ${persona.age} • ${persona.city}`;

    const totalCorpus = persona.deposits.reduce((sum, d) => sum + d.amount, 0);
    const activeCount = persona.deposits.filter(d => d.status === 'active').length;

    document.getElementById('profileStats').innerHTML = `
        <div class="profile-stat">
            <span class="profile-stat-value">${persona.deposits.length}</span>
            <span class="profile-stat-label">Deposits</span>
        </div>
        <div class="profile-stat">
            <span class="profile-stat-value">₹${formatLakhs(totalCorpus)}</span>
            <span class="profile-stat-label">Corpus</span>
        </div>
        <div class="profile-stat">
            <span class="profile-stat-value">${getYearsActive(persona.deposits)}</span>
            <span class="profile-stat-label">Years</span>
        </div>
        <div class="profile-stat">
            <span class="profile-stat-value">${activeCount}</span>
            <span class="profile-stat-label">Active</span>
        </div>
    `;

    document.getElementById('profileStory').innerHTML = `
        <p>${story.narrative}</p>
        <p class="hook">${story.hook}</p>
        <p class="punchline">${story.punchline}</p>
    `;

    showView('profile');
}

async function analyzePersona() {
    const btn = document.querySelector('.btn-analyze');
    btn.textContent = 'Analyzing...';
    btn.disabled = true;

    showView('loading');
    document.getElementById('loadingYears').textContent = getYearsActive(personas.find(p => p.id === currentPersonaId).deposits);

    const steps = [
        'Evaluating consistency',
        'Checking maturity discipline',
        'Measuring growth trajectory',
        'Assessing diversification',
        'Testing tenure intelligence'
    ];

    const stepsContainer = document.getElementById('loadingSteps');
    stepsContainer.innerHTML = '';

    for (let i = 0; i < steps.length; i++) {
        await new Promise(r => setTimeout(r, 500));
        stepsContainer.innerHTML += `<div class="loading-step" id="step${i}">${steps[i]}</div>`;
        document.getElementById(`step${i}`).classList.add('done');
    }

    await new Promise(r => setTimeout(r, 500));

    try {
        const response = await fetch(`/api/personas/${currentPersonaId}/score`);
        currentScore = await response.json();
        showScoreReveal();
    } catch (error) {
        console.error('Failed to get score:', error);
    }

    btn.textContent = 'Analyze FD History';
    btn.disabled = false;
}

function showScoreReveal() {
    showView('scoreReveal');

    document.getElementById('scoreTitle').textContent = `${currentScore.persona.name}'s FD Credit Score`;
    document.getElementById('cibilEquiv').textContent = currentScore.cibil_equivalent;
    document.getElementById('peerPercentile').textContent = `${currentScore.peer_percentile}%`;
    document.getElementById('currentScore').textContent = currentScore.total_score;
    document.getElementById('currentBand').textContent = currentScore.score_band;
    document.getElementById('projectedScore').textContent = currentScore.projected_score;

    animateScoreDial(currentScore.total_score);
    document.getElementById('scoreBand').textContent = currentScore.score_band;

    renderRadarChart(currentScore.components);
    renderComponentCards(currentScore.components);
    renderPatternCard(currentScore.pattern_detected, currentScore.insights);
    renderDepositTimeline();
    renderProducts(currentScore.credit_products);
    renderImprovements(currentScore.improvements);

    setTimeout(() => {
        document.querySelectorAll('.scroll-reveal').forEach(el => {
            el.classList.add('visible');
        });
    }, 100);
}

function animateScoreDial(score) {
    const circumference = 2 * Math.PI * 90;
    const progress = document.querySelector('.dial-progress');
    const offset = circumference - (score / 900) * circumference;

    setTimeout(() => {
        progress.style.strokeDashoffset = offset;
    }, 100);

    let current = 0;
    const target = score;
    const duration = 2000;
    const increment = target / (duration / 16);

    const counter = setInterval(() => {
        current += increment;
        if (current >= target) {
            current = target;
            clearInterval(counter);
        }
        document.getElementById('scoreNumber').textContent = Math.round(current);
    }, 16);

    const scorePercent = (score - 300) / 600;
    document.getElementById('progressFill').style.width = `${scorePercent * 100}%`;
    document.getElementById('potentialMarker').style.left = `${((currentScore.projected_score - 300) / 600) * 100}%`;
}

function renderComponentCards(components) {
    const container = document.getElementById('componentsList');
    const colors = ['#00d4aa', '#0066ff', '#4ecdc4', '#ffe66d', '#ff8a5c'];

    container.innerHTML = components.map((c, i) => `
        <div class="component-card" onclick="this.classList.toggle('expanded')">
            <div class="component-header">
                <span class="component-name">${c.name}</span>
                <span class="component-weight">${Math.round(c.weight * 100)}% weight</span>
            </div>
            <div class="component-bar">
                <div class="component-fill" style="width: ${c.score}%; background: ${colors[i]}"></div>
            </div>
            <span class="component-score">${c.score}</span>
            <div class="component-details">
                ${c.sub_metrics.map(sm => `
                    <div class="sub-metric">
                        <span>${sm.name}</span>
                        <span class="sub-metric-score">${sm.score}/${sm.max}</span>
                    </div>
                `).join('')}
            </div>
        </div>
    `).join('');
}

function renderPatternCard(pattern, insights) {
    const card = document.getElementById('patternCard');
    const patternClass = pattern.toLowerCase().replace(/_/g, '-');
    const narrative = patternNarratives[pattern] || insights[0];

    card.className = `pattern-card ${patternClass}`;
    card.innerHTML = `
        <div class="pattern-title">${pattern.replace(/_/g, ' ')}</div>
        <div class="pattern-narrative">${narrative}</div>
    `;
}

function renderDepositTimeline() {
    const ctx = document.getElementById('timelineChart').getContext('2d');
    const deposits = currentScore.persona.deposit_count > 0 
        ? personas.find(p => p.id === currentPersonaId)?.deposits || []
        : [];

    const yearlyData = {};
    deposits.forEach(d => {
        const year = d.start_date.substring(0, 4);
        if (!yearlyData[year]) yearlyData[year] = { matured: 0, active: 0, broken: 0 };
        
        if (d.status === 'matured') yearlyData[year].matured += d.amount;
        else if (d.status === 'active') yearlyData[year].active += d.amount;
        else if (d.withdrawn_date) yearlyData[year].broken += d.amount;
    });

    const years = Object.keys(yearlyData).sort();
    const matured = years.map(y => yearlyData[y].matured / 100000);
    const active = years.map(y => yearlyData[y].active / 100000);
    const broken = years.map(y => yearlyData[y].broken / 100000);

    if (timelineChart) timelineChart.destroy();

    timelineChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: years,
            datasets: [
                { label: 'Matured', data: matured, backgroundColor: '#00d4aa' },
                { label: 'Active', data: active, backgroundColor: '#0066ff' },
                { label: 'Broken', data: broken, backgroundColor: '#ff6b6b' }
            ]
        },
        options: {
            responsive: true,
            scales: {
                x: { stacked: true, grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#9ca3af' } },
                y: { stacked: true, grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#9ca3af', callback: v => '₹' + v + 'L' } }
            },
            plugins: { legend: { display: false } }
        }
    });
}

function renderProducts(products) {
    const grid = document.getElementById('productsGrid');
    const icons = { 'Secured': 'S', 'Credit Line': 'C', 'Unsecured': 'U' };

    grid.innerHTML = products.map(p => `
        <div class="product-card">
            <div class="product-icon">${icons[p.type] || 'P'}</div>
            <div class="product-name">${p.name}</div>
            <div class="product-type">${p.type}</div>
            <div class="product-details">
                Limit: ₹${p.limit.toLocaleString()}<br>
                Interest: ${p.interest}%
            </div>
            <span class="product-badge eligible">✅ Eligible</span>
        </div>
    `).join('');
}

function renderImprovements(improvements) {
    const list = document.getElementById('improvementsList');
    list.innerHTML = improvements.map(i => `
        <div class="improvement-item">
            <span class="improvement-action">${i.action}</span>
            <div class="improvement-meta">
                <span class="difficulty-badge ${i.difficulty.toLowerCase()}">${i.difficulty}</span>
                <span class="improvement-points">+${i.points_delta} pts</span>
            </div>
        </div>
    `).join('');
}

function formatLakhs(amount) {
    return (amount / 100000).toFixed(1) + 'L';
}

function showView(viewId) {
    if (viewId === 'scoreReveal') {
        resetScoreReveal();
    } else if (viewId === 'loading') {
        resetLoading();
    }
    
    document.querySelectorAll('.view').forEach(v => {
        v.classList.add('hidden');
        v.classList.remove('active');
    });
    const view = document.getElementById(viewId);
    view.classList.remove('hidden');
    view.classList.add('active');
    
    if (viewId === 'landing') {
        renderPersonas();
        if (lastScrollPosition > 0) {
            setTimeout(() => window.scrollTo(0, lastScrollPosition), 50);
        }
    } else if (viewId === 'profile') {
        window.scrollTo(0, 0);
    }
    
    window.scrollTo(0, 0);
}

function resetScoreReveal() {
    document.getElementById('scoreNumber').textContent = '0';
    document.getElementById('scoreBand').textContent = '--';
    document.getElementById('cibilEquiv').textContent = '--';
    document.getElementById('peerPercentile').textContent = '--';
    document.getElementById('currentScore').textContent = '--';
    document.getElementById('currentBand').textContent = '--';
    document.getElementById('projectedScore').textContent = '--';
    document.getElementById('scoreTitle').textContent = '';
    
    const dialProgress = document.querySelector('.dial-progress');
    if (dialProgress) dialProgress.style.strokeDashoffset = 565.48;
    
    const progressFill = document.getElementById('progressFill');
    if (progressFill) progressFill.style.width = '0';
    
    const potentialMarker = document.getElementById('potentialMarker');
    if (potentialMarker) potentialMarker.style.left = '0';
    
    const componentsList = document.getElementById('componentsList');
    if (componentsList) componentsList.innerHTML = '';
    
    const patternCard = document.getElementById('patternCard');
    if (patternCard) patternCard.innerHTML = '';
    
    const productsGrid = document.getElementById('productsGrid');
    if (productsGrid) productsGrid.innerHTML = '';
    
    const improvementsList = document.getElementById('improvementsList');
    if (improvementsList) improvementsList.innerHTML = '';
    
    if (radarChart) { 
        radarChart.destroy(); 
        radarChart = null; 
    }
    if (timelineChart) { 
        timelineChart.destroy(); 
        timelineChart = null; 
    }
    
    document.querySelectorAll('.scroll-reveal').forEach(el => el.classList.remove('visible'));
}

function resetProfilePage() {
    document.getElementById('profileAvatar').innerHTML = '';
    document.getElementById('profileName').textContent = '';
    document.getElementById('profileMeta').textContent = '';
    document.getElementById('profileStats').innerHTML = '';
    document.getElementById('profileStory').innerHTML = '';
    const btn = document.querySelector('.btn-analyze');
    if (btn) {
        btn.textContent = 'Analyze FD History';
        btn.disabled = false;
    }
}

function resetLoading() {
    const stepsContainer = document.getElementById('loadingSteps');
    if (stepsContainer) stepsContainer.innerHTML = '';
    const loadingYears = document.getElementById('loadingYears');
    if (loadingYears) loadingYears.textContent = '';
}

function resetFormWizard() {
    document.getElementById('inputName').value = '';
    document.getElementById('inputAge').value = '';
    document.getElementById('inputCity').value = '';
    customDeposits = [];
    wizardStep = 1;
    document.querySelectorAll('.wizard-step-content').forEach(el => el.classList.add('hidden'));
    document.getElementById('step1').classList.remove('hidden');
    document.getElementById('wizardStep').textContent = 'Step 1 of 3';
    document.getElementById('wizardProgressFill').style.width = '33%';
    document.getElementById('depositsList').innerHTML = '';
}

function scrollToSection(sectionId) {
    document.getElementById(sectionId).scrollIntoView({ behavior: 'smooth' });
}

// Wizard functions
document.getElementById('depositStatus').addEventListener('change', function() {
    document.getElementById('withdrawnDateGroup').classList.toggle('hidden', this.value !== 'withdrawn');
});

function nextStep(step) {
    if (step === 2) {
        const name = document.getElementById('inputName').value;
        const age = document.getElementById('inputAge').value;
        const city = document.getElementById('inputCity').value;
        if (!name || !age || !city) {
            alert('Please fill in all fields');
            return;
        }
    }

    if (step === 3 && customDeposits.length === 0) {
        alert('Please add at least one deposit');
        return;
    }

    wizardStep = step;
    document.querySelectorAll('.wizard-step-content').forEach(el => el.classList.add('hidden'));
    document.getElementById(`step${step}`).classList.remove('hidden');
    document.getElementById('wizardStep').textContent = `Step ${step} of 3`;
    document.getElementById('wizardProgressFill').style.width = `${(step / 3) * 100}%`;

    if (step === 3) {
        renderReview();
    }
}

function addDeposit() {
    const deposit = {
        type: document.getElementById('depositType').value,
        bank: document.getElementById('depositBank').value,
        amount: parseFloat(document.getElementById('depositAmount').value),
        tenure_months: parseInt(document.getElementById('depositTenure').value),
        interest_rate: parseFloat(document.getElementById('depositRate').value),
        start_date: document.getElementById('depositStartDate').value,
        maturity_date: getMaturityDate(document.getElementById('depositStartDate').value, parseInt(document.getElementById('depositTenure').value)),
        status: document.getElementById('depositStatus').value,
        withdrawn_date: document.getElementById('depositStatus').value === 'withdrawn' ? document.getElementById('withdrawnDate').value : null
    };

    if (!deposit.amount || !deposit.start_date) {
        alert('Please fill in all required fields');
        return;
    }

    customDeposits.push(deposit);
    renderDepositsList();
    clearDepositForm();
}

function getMaturityDate(startDate, tenure) {
    if (!startDate) return '';
    const start = new Date(startDate);
    start.setMonth(start.getMonth() + tenure);
    return start.toISOString().split('T')[0];
}

function renderDepositsList() {
    const list = document.getElementById('depositsList');
    list.innerHTML = customDeposits.map((d, i) => `
        <div class="deposit-item">
            <div class="deposit-info">
                <span class="deposit-amount">₹${d.amount.toLocaleString()} ${d.type}</span>
                <span class="deposit-meta"> at ${d.bank} • ${d.tenure_months}mo • ${d.status}</span>
            </div>
            <button onclick="removeDeposit(${i})" style="background:none;border:none;color:#ff6b6b;cursor:pointer">✕</button>
        </div>
    `).join('');
}

function removeDeposit(index) {
    customDeposits.splice(index, 1);
    renderDepositsList();
}

function clearDepositForm() {
    document.getElementById('depositAmount').value = '';
    document.getElementById('depositRate').value = '';
    document.getElementById('depositStartDate').value = '';
}

function renderReview() {
    const name = document.getElementById('inputName').value;
    const city = document.getElementById('inputCity').value;
    document.getElementById('reviewSummary').innerHTML = `
        <p><strong>You:</strong> ${name}, ${city}</p>
        <p><strong>Deposits:</strong> ${customDeposits.length}</p>
        ${customDeposits.map(d => `<p>• ₹${d.amount.toLocaleString()} ${d.type} at ${d.bank} (${d.status})</p>`).join('')}
    `;
}

async function calculateCustomScore() {
    showView('loading');
    document.getElementById('loadingYears').textContent = getYearsActive(customDeposits);

    const steps = ['Evaluating consistency', 'Checking maturity discipline', 'Measuring growth trajectory', 'Assessing diversification', 'Testing tenure intelligence'];
    const stepsContainer = document.getElementById('loadingSteps');
    stepsContainer.innerHTML = '';

    for (let i = 0; i < steps.length; i++) {
        await new Promise(r => setTimeout(r, 500));
        stepsContainer.innerHTML += `<div class="loading-step" id="step${i}">${steps[i]}</div>`;
        document.getElementById(`step${i}`).classList.add('done');
    }

    await new Promise(r => setTimeout(r, 500));

    const history = {
        user_id: 'custom',
        name: document.getElementById('inputName').value,
        age: parseInt(document.getElementById('inputAge').value),
        city: document.getElementById('inputCity').value,
        deposits: customDeposits
    };

    try {
        const response = await fetch('/api/score', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(history)
        });
        currentScore = await response.json();
        currentScore.persona = {
            name: history.name,
            age: history.age,
            occupation: 'Custom User',
            city: history.city,
            deposit_count: customDeposits.length,
            total_corpus: customDeposits.reduce((s, d) => s + d.amount, 0),
            years_active: getYearsActive(customDeposits),
            active_deposits: customDeposits.filter(d => d.status === 'active').length
        };
        showScoreReveal();
    } catch (error) {
        console.error('Failed to calculate score:', error);
    }
}

function setupScrollObserver() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
            }
        });
    }, { threshold: 0.1 });

    document.querySelectorAll('.scroll-reveal').forEach(el => observer.observe(el));
}

init();