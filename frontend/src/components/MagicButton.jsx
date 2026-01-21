import React from 'react';
import { Sparkles } from 'lucide-react';
import { motion } from 'framer-motion';

const MagicButton = ({ onFill, label = "Авто-заполнение" }) => {
    return (
        <motion.button
            type="button"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={onFill}
            className="btn"
            style={{
                background: 'linear-gradient(135deg, #FFD700 0%, #FFA500 100%)',
                color: '#000',
                fontWeight: 'bold',
                border: 'none'
            }}
        >
            <Sparkles size={18} />
            {label}
        </motion.button>
    );
};

export default MagicButton;
