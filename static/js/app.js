document.addEventListener('DOMContentLoaded', function() {
    // 表单提交处理
    document.querySelectorAll('form').forEach(form => {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = new FormData(this);
            const action = this.getAttribute('action');
            const method = this.getAttribute('method');
            
            try {
                const response = await fetch(action, {
                    method: method,
                    body: formData
                });
                
                if (response.redirected) {
                    window.location.href = response.url;
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
    });

    // 任务完成状态切换动画
    document.querySelectorAll('.task').forEach(task => {
        task.addEventListener('click', function(e) {
            if (e.target.tagName === 'BUTTON') {
                this.style.opacity = '0.5';
                setTimeout(() => {
                    this.style.opacity = '1';
                }, 300);
            }
        });
    });

    // 输入框验证
    const taskInput = document.querySelector('input[name="content"]');
    if (taskInput) {
        taskInput.addEventListener('input', function() {
            if (this.value.trim().length > 0) {
                this.setCustomValidity('');
            } else {
                this.setCustomValidity('请输入任务内容');
            }
        });
    }
});
