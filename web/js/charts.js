function renderRadarChart(components) {
    const ctx = document.getElementById('radarChart').getContext('2d');
    
    const labels = components.map(c => c.name);
    const scores = components.map(c => c.score);
    
    if (radarChart) {
        radarChart.destroy();
    }
    
    radarChart = new Chart(ctx, {
        type: 'radar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Score',
                data: scores,
                backgroundColor: 'rgba(0, 212, 170, 0.2)',
                borderColor: 'rgba(0, 212, 170, 1)',
                borderWidth: 2,
                pointBackgroundColor: 'rgba(0, 212, 170, 1)',
                pointBorderColor: '#fff',
                pointHoverBackgroundColor: '#fff',
                pointHoverBorderColor: 'rgba(0, 212, 170, 1)'
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            scales: {
                r: {
                    beginAtZero: true,
                    max: 100,
                    ticks: {
                        stepSize: 20,
                        color: 'rgba(255, 255, 255, 0.5)',
                        backdropColor: 'transparent'
                    },
                    pointLabels: {
                        color: 'rgba(255, 255, 255, 0.8)',
                        font: {
                            size: 11,
                            family: 'Inter'
                        }
                    },
                    grid: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    },
                    angleLines: {
                        color: 'rgba(255, 255, 255, 0.1)'
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            },
            animation: {
                duration: 1500,
                easing: 'easeOutQuart'
            }
        }
    });
}