import {
  ConfigProvider,
  Button,
  Form,
  Input,
  message,
  Layout,
  Menu,
  Row,
  Col,
  Table,
  Card,
  Pagination,
  Modal,
  Select,
  Switch,
  Upload
} from 'ant-design-vue'

message.config({
  top: '60px',
  duration: 2,
  maxCount: 3
})

export default function (app) {
  app.prototype.$message = message
  app.prototype.$confirm = Modal.confirm // 对话框

  app.use(Button)
  app.use(Form)
  app.use(Input)
  app.use(Layout)
  app.use(Menu)
  app.use(Row)
  app.use(Col)
  app.use(Table)
  app.use(Card)
  app.use(Pagination)
  app.use(ConfigProvider)
  app.use(Modal)
  app.use(Select)
  app.use(Switch)
  app.use(Upload)
}


