import {useState,useEffect} from 'react';
import { Avatar, Result, Popover, Card, Descriptions, Form, Switch, Input, Tooltip, Radio  } from 'antd';
import { MessageFilled, SettingOutlined, WomanOutlined , ManOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { EventsOn, EventsOff, EventsEmit } from '../../wailsjs/runtime/runtime'
import { Logger } from "../../wailsjs/go/main/App";
import './Home.css'
const { Meta } = Card;
const { TextArea } = Input;

function Home() {
  const [settingForm] = Form.useForm();
  const [silderType, setSilderType] = useState<string>('setting')
  const [pageStatus, setPageStatus] = useState<number>(0)
  const [unid, setUnid] = useState<number>(123)
  const [sex, setSex] = useState<number>(0)
  const [nickName, setNickName ] = useState<string>('')
  const [friendList, setFriendList] = useState<any[]>([])
  const [activeFriend, setActiveFriend] = useState<number>(-1)
  const [firiendItem, setFriendItem] = useState<any>(null)
  const [settingShow, setSettingShow] = useState<boolean>(false)
  const [settingConfig, setSettingConfig] = useState<any>({ 
    auto_reply: false, 
    auto_reply_group: false,
    auto_bot: 'nobot',
    auto_desc: '正在忙ing',
  })
  useEffect(() => {
    const getUserInfo = (userInfo: any) => {
      Logger(userInfo)
      if (userInfo.StatusCode === 200) {
        setPageStatus(0)
        setSex(userInfo.Info.Sex)
        setNickName(userInfo.Info.NickName)
        setFriendList(userInfo.Firends)
        setUnid(userInfo.Unid)
        EventsEmit('config:ready', userInfo.Unid)
      } else {
        setPageStatus(-1)
      }
    }
    EventsOn('info:get', getUserInfo)
    EventsEmit('info:ready')
    return ()=> {
      EventsOff('info:get', 'getUserInfo')
    }
  }, [])

  useEffect(() => {
    const getConfig = (config: any) => {
      Logger(config)
      if (typeof config === 'object' && Object.keys(config).length) {
        setSettingConfig(config)
        setSettingShow(config.auto_reply)
      }
    }
    EventsOn('config:get', getConfig)
    return ()=> {
      EventsOff('config:get', 'getUserInfo')
    }
  }, [])
  
  const handleSilder = (key: string) => {
    if (pageStatus === -1) return
    setSilderType(key)
    setActiveFriend(-1)
    setFriendItem(null)
  }
  const Error = () => {
    return <div className='flex-1'>
      <Result
          status="500"
          title="500"
          subTitle="服务异常"
        />
    </div>
  }

  const handleFriendItem = (index: number) => {
    setActiveFriend(index)
    setFriendItem(friendList[index])
  }
  const Friend = () => {
    return <>
      <div className="friend-content">
        { friendList.map((item: any, index: number) => {
          return <div>
            <Card  hoverable={true} onClick={ () =>  handleFriendItem(index) } className={ index === activeFriend ? 'meta-active meta-content' : 'meta-content' } >
              <Meta
                avatar={<Avatar icon={SexNode(item.Sex)}/>}
                title={ item.RemarkName ? item.RemarkName : item.NickName}
                description={ '昵称:' +  item.NickName}
              />
            </Card>
          </div>
        })}
      </div>
      { firiendItem ? <div  className="setting-content">
        <Descriptions column={1} >
          <Descriptions.Item label="昵称">{ firiendItem.NickName }</Descriptions.Item>
          <Descriptions.Item label="备注">{ firiendItem.RemarkName }</Descriptions.Item>
          <Descriptions.Item label="签名">{ firiendItem.Signature }</Descriptions.Item>
        </Descriptions>
        {/* <Form labelCol={{ span: 4 }}  layout="horizontal">
          <Form.Item label="自动回复" valuePropName="checked">
            <Switch />
          </Form.Item>
          <div>
            <Form.Item label={ <Tooltip title="低于机器人回复优先级,100字以内">
                <span>自动回复文案</span>
              </Tooltip> }>
              <TextArea rows={4} maxLength={100}/>
            </Form.Item>
            <Form.Item label="机器人">
              <Radio.Group>
                <Radio value="nobot"> 无 </Radio>
                <Radio value="tuling"> 图灵 </Radio>
                <Radio value="chatgpt"> ChatGPT </Radio>
              </Radio.Group>
            </Form.Item>
          </div>
        </Form> */}
      </div> : '' }
 
    </>
  }
  const onSettingFinish = (values: any) => {
    const newValues = { unid, ...values }
    Logger(newValues)
    EventsEmit('config:save', newValues)
  }

  const onChange = (check: boolean) => {
    setSettingShow(check)
    settingForm.submit()
  }
  const Setting = () => {
    return  <div className="setting-content">
        <Form form={settingForm} labelCol={{ span: 4 }} initialValues={ settingConfig } layout="horizontal" onFinish={onSettingFinish}>
          <Form.Item label="自动回复" name="auto_reply" valuePropName="checked">
            <Switch  onChange={ (checked)=>  onChange(checked)}/>
          </Form.Item>
          <div style={ { visibility: settingShow ? 'visible' : 'hidden'}}>
            <Form.Item label="群@自动回复" name="auto_reply_group" valuePropName="checked">
              <Switch  onChange= { ()=> settingForm.submit() } />
            </Form.Item>
            <Form.Item name="auto_desc" label={ <Tooltip title="低于机器人回复优先级,100字以内">
                <span>自动回复文案<InfoCircleOutlined style={{ color: '#66666 !important',}} /></span>
              </Tooltip> }>
              <TextArea rows={4} maxLength={100} onBlur={ ()=> settingForm.submit() }/>
            </Form.Item>
            <Form.Item name="auto_bot" label="机器人">
              <Radio.Group onChange= { ()=> settingForm.submit() }>
                <Radio value="nobot"> 无 </Radio>
                <Radio value="tuling"> 图灵 </Radio>
                <Radio disabled value="chatgpt"> ChatGPT </Radio>
              </Radio.Group>
            </Form.Item>
            <Form.Item name="tuling_api_key" label={ <Tooltip title="http://www.turingapi.com/">
                <span>图灵机器人APP_KEY <InfoCircleOutlined style={{ color: '#66666 !important',}} /></span>
              </Tooltip> }>
            <Input placeholder="填写app_key,使用图灵机器人"  onBlur={ ()=> settingForm.submit() } />
            </Form.Item>
          </div>
        </Form>
    </div>
  }

  const Content = ()=> {
    return silderType === 'friend' ? Friend() : Setting();
  }

  const SexNode = (sex: number) => {
    return sex === 1 ? <ManOutlined /> : <WomanOutlined />
  }
  const WxHeader = () => {
    return nickName !== '' ? (<Popover content={nickName} placement="rightTop" trigger="hover">
        <Avatar size={48} icon={ SexNode(sex) }></Avatar>
      </Popover>)
    : <Avatar size={48} icon={ SexNode(sex) }></Avatar>
  }
  return <>
    <div className="layout">
      <div className="silder">
        { WxHeader() }
        {/* <MessageFilled onClick={()=> handleSilder('friend')} className={ `icon ${silderType === 'friend' && 'icon-active'}` } /> */}
        <SettingOutlined onClick={()=> handleSilder('setting')} className={ `icon ${silderType === 'setting' && 'icon-active'}` }/>
      </div>
      { pageStatus === -1 ? Error() : Content()}

    </div>
   
  </>
}
export default Home
