import { AxiosError } from 'axios'
import { useAPI, ApiResponse, apiUrl } from '../../api/api'

const {onGetConfigHash,onBlockAction,onNewRedirectAction,onDeleteRedirectAction,onReloadActions,onUnblockAction} = useAPI()

export {onBlockAction,onDeleteRedirectAction,onGetConfigHash,onNewRedirectAction,onReloadActions,onUnblockAction}