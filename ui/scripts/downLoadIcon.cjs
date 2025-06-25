const fs = require('fs')
const path = require('path')

const Axios = require('axios')

async function downloadFile(url) {
  try {
    const iconPath = path.resolve(__dirname, '../src/assets/fonts/iconfont.js')
    const iconDir = path.dirname(iconPath)
    
    // 检查目录是否存在，如果不存在则创建
    if (!fs.existsSync(iconDir)) {
      console.log(`目录 ${iconDir} 不存在，正在创建...`)
      fs.mkdirSync(iconDir, { recursive: true })
      console.log('目录创建成功')
    }
    
    console.log(`开始下载图标文件到: ${iconPath}`)
    
    const writer = fs.createWriteStream(iconPath)
    const response = await Axios({
      url: `https:${url}`,
      method: 'GET',
      responseType: 'stream',
      timeout: 30000, // 30秒超时
    })
    
    response.data.pipe(writer)
    
    return new Promise((resolve, reject) => {
      writer.on('finish', () => {
        console.log('图标文件下载成功！')
        resolve()
      })
      writer.on('error', (err) => {
        console.error('写入文件时出错:', err.message)
        reject(err)
      })
    })
  } catch (error) {
    console.error('下载过程中出错:', error.message)
    throw error
  }
}

async function main() {
  const argument = process.argv.splice(2)
  
  if (!argument[0]) {
    console.error('错误: 请提供下载URL作为参数')
    console.log('使用方法: node downLoadIcon.cjs <url>')
    process.exit(1)
  }
  
  try {
    await downloadFile(argument[0])
    console.log('所有操作完成！')
  } catch (error) {
    console.error('脚本执行失败:', error.message)
    process.exit(1)
  }
}

main()
